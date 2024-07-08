from pymilvus import db, MilvusException, connections, utility, DataType, FieldSchema, CollectionSchema, Collection
from angle_emb import Prompts
import os


def createMilvusDatabase(databaseName):
    connections.connect(host="localhost", port="19530")
    try:
        db.using_database(databaseName)
    except MilvusException:
        db.create_database(databaseName)
        msg = f"Database '{databaseName}' created."
    else:
        msg = f"Database '{databaseName}' already exists."
    connections.disconnect(alias="default")
    return msg
    
    
def createMilvusCollection(databaseName, collectionName, dim=1024):
    connections.connect(host="localhost", port="19530")
    db.using_database(databaseName)
    
    if utility.has_collection(collectionName):
        connections.disconnect(alias="default")
        return False
    
    id_field = FieldSchema(
        name="id", dtype=DataType.INT64, is_primary=True, auto_id=True
    )
    vector_field = FieldSchema(name="vector", dtype=DataType.FLOAT_VECTOR, dim=dim)
    page_id_field = FieldSchema(name="page_id", dtype=DataType.INT64)

    schema = CollectionSchema(
        fields=[id_field, vector_field, page_id_field],
        description=f"Novel collection for novel ID {collectionName}",
    )
    Collection(name=collectionName, schema=schema)
    connections.disconnect(alias="default")
    return True


def setUpMilvusCollection(databaseName, collectionName, bookID, embeddingModel):
    base_dir = os.path.abspath(f"../webGateway/cache/{bookID}/pages")
    pageLst = []
    pageEbdLst = []
    for filename in os.listdir(base_dir):
        if filename.endswith(".txt"):
            page_number = int(filename.split('_')[1].split('.')[0])
            
            file_path = os.path.join(base_dir, filename)
            with open(file_path, 'r') as file:
                print(f"Processing page {page_number}...")
                file.read()
                pageLst.append(page_number)
                pageEbdLst.append(embeddingModel.encode([file.read()])[0])
    addChunksToCollection(databaseName, collectionName, pageEbdLst, pageLst)


def addChunksToCollection(databaseName, collectionName, vectors, pageLst):
    connections.connect(host="localhost", port="19530")
    db.using_database(databaseName)
    collection = Collection(name=collectionName)
    data = [vectors, pageLst]
    collection.insert(data)
    index_params = {
        "metric_type": "IP",
        "index_type": "IVF_FLAT",
        "params": {"nlist": 16384}
    }
    
    collection.create_index(field_name="vector", index_params=index_params)
    
    collection.release()
    connections.disconnect(alias="default")
    
def retriveKnowledge(databaseName, collectionName, bookID, query, curPage, ebdModel):
    results = searchTopChunks(databaseName, collectionName, query, curPage, ebdModel)
    return retrieveText(bookID, results)

def searchTopChunks(databaseName, collectionName, query, curPage, ebdModel, topk=5):
    connections.connect(host="localhost", port="19530")
    db.using_database(databaseName)
    collection = Collection(name=collectionName)
    collection.load()
    qv = ebdModel.encode(Prompts.C.format(text=query))
    search_params = {"metric_type": "IP"}
    expr = f"page_id <= {curPage}"
    results = collection.search(
        data= qv,
        anns_field="vector",
        param=search_params,
        limit=topk,
        expr=expr,
        consistency_level="Strong",
    )
    id_list = results[0].ids
    query_expression = f"id in {id_list}"
    results = collection.query(
        expr=query_expression, output_fields=["page_id"]
    )
    collection.release()
    connections.disconnect(alias="default")
    return results


def retrieveText(bookID, resultList):
    textLst = []
    for pageNum in resultList:
        file_path = os.path.abspath(f"../webGateway/cache/{bookID}/pages/page_{pageNum['page_id']}.txt")
        with open(file_path, 'r') as file:
            textLst.append(file.read())
    return textLst