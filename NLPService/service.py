from flask import Flask, request, jsonify
from angle_emb import AnglE
from utils import milvus_helper
import logging

app = Flask(__name__)
logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s %(levelname)s: %(message)s",
    datefmt="%Y-%m-%d %H:%M:%S",
)
app.logger.addHandler(logging.StreamHandler())
ebdModel = AnglE.from_pretrained('WhereIsAI/UAE-Large-V1', pooling_strategy='cls').cuda()
@app.route('/chat', methods=['POST'])
def chat():
    data = request.get_json() 
    model = data.get('model')
    apikey = data.get('apikey')
    bookID = data.get('bookID')
    progress = data.get('progress')
    userMessage = data.get('userMessage')
    topText = milvus_helper.retriveKnowledge("reading_copilot", f"book_{bookID}", bookID, userMessage, progress, ebdModel)
    responseMessage = f"book: {bookID}, page: {progress}, model: {model}, repeat: {userMessage}, most similar text: {topText}"
    return jsonify({"responseMessage": responseMessage, "error": ""})

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