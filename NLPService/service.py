from flask import Flask, request, jsonify
from angle_emb import AnglE
from utils import milvus_helper, const
from vllm import LLM, SamplingParams
import logging, os


app = Flask(__name__)
logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s %(levelname)s: %(message)s",
    datefmt="%Y-%m-%d %H:%M:%S",
)
app.logger.addHandler(logging.StreamHandler())
ebdModel = AnglE.from_pretrained('WhereIsAI/UAE-Large-V1', pooling_strategy='cls').cuda()

localModelPath = os.getenv('LOCAL_LLM_PATH', '')
if localModelPath != '':
    localmodel = LLM(model=localModelPath)
sampling_params = SamplingParams(temperature=0.7, top_p=0.9, max_tokens=800)

@app.route('/chat', methods=['POST'])
def chat():
    data = request.get_json() 
    model = data.get('model')
    apikey = data.get('apikey')
    bookID = data.get('bookID')
    progress = data.get('progress')
    userMessage = data.get('userMessage')
    topText = milvus_helper.retriveKnowledge("reading_copilot", f"book_{bookID}", bookID, userMessage, progress, ebdModel)
    user_message = const.COPILOT.format(
        contents="\n".join([page for page in topText]), question=userMessage
    )
    prompt = const.TEMPLATE.format(user_message=user_message)
    outputs = model.generate(prompt, sampling_params)[0].outputs[0].text
    return jsonify({"responseMessage": outputs, "error": ""})

@app.route('/build', methods=['POST'])
def milvus():
    data = request.get_json()
    bookID = data['bookID']
    collectionName = f"book_{bookID}"
    app.logger.info(f"Building collection {collectionName}...")
    needSetUp = milvus_helper.createMilvusCollection("reading_copilot", collectionName)
    app.logger.info(f"Collection {collectionName} has been created.")
    if needSetUp:
        app.logger.info(f"Setting up collection {collectionName}...")
        milvus_helper.setUpMilvusCollection("reading_copilot", collectionName, bookID, ebdModel)
        app.logger.info(f"Collection {collectionName} has been set up.")
        response = jsonify({"msg": f"Collection {collectionName} has been set up."})
    else:
        app.logger.info(f"Collection {collectionName} already exists.")
        response = jsonify({"msg": f"Collection {collectionName} already exists."})
    return response

if __name__ == '__main__':    
    context = ('tls/cert.pem', 'tls/key.pem')
    app.logger.info("Starting NLP service...")
    app.logger.info("Connecting to Milvus...")
    app.logger.info(milvus_helper.createMilvusDatabase("reading_copilot"))
    app.run(port=4010, ssl_context=context)