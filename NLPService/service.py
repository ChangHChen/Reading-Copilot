from flask import Flask, request, jsonify

app = Flask(__name__)

@app.route('/chat', methods=['POST'])
def chat():
    data = request.get_json() 
    model = data.get('model')
    apikey = data.get('apikey')
    bookID = data.get('bookID')
    progress = data.get('progress')
    userMessage = data.get('userMessage')
    
    responseMessage = f"book: {bookID}, page: {progress}, model: {model}, repeat: {userMessage}"
    return jsonify({"responseMessage": responseMessage, "error": ""})

@app.route('/milvus', methods=['POST'])
def milvus():
    data = request.get_json()
    bookID = data['bookID']
    
    return jsonify({"err": ""})

if __name__ == '__main__':
    context = ('tls/cert.pem', 'tls/key.pem')
    app.run(port=4010, ssl_context=context)