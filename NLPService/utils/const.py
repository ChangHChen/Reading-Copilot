COPLIOT = (
        "Start of Contents:\n{contents}\nEnd of Contents\n"
        "You will be given some segments of a novel followed by a question asked by the user.\n"
        "Your task is to answer as faithfully as you can.\n"
        "Instructions:\n"
        "- Step 1: Carefully read and understand the provided segment from the novel.\n"
        "- Step 2: Identify the parts of the segment that are related to the question asked.\n"
        "- Step 3: Based on the information in the content, formulate a detailed and structured answer to the question.\n"
        "- Step 4: If the question cannot be answered using the content provided, respond with 'I cannot answer the question based on the information provided.'\n"
        "Question:\n"
        "{question}\n"
    )

TEMPLATE = "<|im_start|>system\nYou are an AI assistant that follows instructions very well. Finish the tasks that given by the user as faithfully as you can.<|im_end|>\n<|im_start|>user\n{user_message}<|im_end|>\n<|im_start|>assistant\n"