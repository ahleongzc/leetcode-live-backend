import re
import random
import inflect

input_file = "./data.txt"
output_file = "../internal/repo/fasttext/labels.txt"

p = inflect.engine()

with open(input_file, "r", encoding="utf-8") as file:
    lines = file.readlines()

cleaned_lines = []
for line in lines:
    cleaned_line = re.sub(r"[^\w\s']", " ", line).lower()
    cleaned_line = re.sub(r"\b\d+\b", lambda x: p.number_to_words(x.group()), cleaned_line)
    cleaned_lines.append(cleaned_line)

random.shuffle(cleaned_lines)

with open(output_file, "w", encoding="utf-8") as file:
    file.writelines(cleaned_lines)

print("Cleaned, converted numbers to words, and shuffled training data")