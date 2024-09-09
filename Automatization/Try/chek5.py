import os
from git import Repo
import javalang
import re

def analyze_exception_handling(project_path):
    results = []

    for root, dirs, files in os.walk(project_path):
        for file in files:
            if file.endswith('.java'):
                file_path = os.path.join(root, file)
                try_catch_stats = analyze_file(file_path)

                results.append({
                    "metric_id": "try_catch_usage",
                    "git_author": os.path.splitext(file)[0],
                    "score": try_catch_stats["score"],
                    "evidence": try_catch_stats["evidence"]
                })

    return results

def analyze_file(file_path):
    try_catch_evidence = {
        "correctly_used": 0,
        "incorrectly_used": 0,
        "missing_use": 0,
        "evidence": []
    }

    with open(file_path, 'r') as file:
        file_content = file.read()

    try:
        tree = javalang.parse.parse(file_content)
        for path, node in tree:
            if isinstance(node, javalang.tree.TryStatement):
                if node.catches:
                    try_catch_evidence["correctly_used"] += 1
                    try_catch_evidence["evidence"].append({
                        "commit_id": None,
                        "file": os.path.basename(file_path),
                        "line": tree.lineno(node)
                    })
                else:
                    try_catch_evidence["incorrectly_used"] += 1
                    try_catch_evidence["evidence"].append({
                        "commit_id": None,
                        "file": os.path.basename(file_path),
                        "line": tree.lineno(node)
                    })
            elif isinstance(node, javalang.tree.MethodInvocation) or isinstance(node, javalang.tree.BinaryOperation) or isinstance(node, javalang.tree.ArraySelector) or isinstance(node, javalang.tree.Cast):
                try_catch_evidence["missing_use"] += 1
                try_catch_evidence["evidence"].append({
                    "commit_id": None,
                    "file": os.path.basename(file_path),
                    "line": tree.lineno(node)
                })
    except javalang.parser.JavaSyntaxError:
        pass

    total_usage = try_catch_evidence["correctly_used"] + try_catch_evidence["incorrectly_used"] + try_catch_evidence["missing_use"]
    if total_usage == 0:
        try_catch_evidence["score"] = "N/A"
    elif try_catch_evidence["correctly_used"] == total_usage:
        try_catch_evidence["score"] = "5"
    elif try_catch_evidence["correctly_used"] == total_usage - try_catch_evidence["missing_use"]:
        try_catch_evidence["score"] = "4"
    elif try_catch_evidence["incorrectly_used"] == 0:
        try_catch_evidence["score"] = "3"
    elif try_catch_evidence["correctly_used"] == 0:
        try_catch_evidence["score"] = "2"
    else:
        try_catch_evidence["score"] = "1"

    return try_catch_evidence

def main():
    project_path = '/Users/damianmarquez/Documents/up/fury_ads-search-charge-bonification/src/main'
    results = analyze_exception_handling(project_path)
    print(results)

if __name__ == '__main__':
    main()