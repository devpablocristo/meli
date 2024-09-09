import os
from git import Repo
import javalang

def analyze_try_catch_usage(repo_path):
    repo = Repo(repo_path)
    committer_scores = {}

    for commit in repo.iter_commits():
        print(1)
        committer_name = commit.committer.name
        if committer_name not in committer_scores:
            committer_scores[committer_name] = {
                'try_catch_usage': {
                    'score': 'N/A',
                    'correctly_used': 0,
                    'incorrectly_used': 0,
                    'missing_use': 0
                }
            }

        for diff in commit.diff(None):
            if diff.b_path.endswith('.java'):
                print(2)
                try_catch_stats = analyze_file(os.path.join(repo_path, diff.b_path))
                committer_scores[committer_name]['try_catch_usage']['correctly_used'] += try_catch_stats['correctly_used']
                committer_scores[committer_name]['try_catch_usage']['incorrectly_used'] += try_catch_stats['incorrectly_used']
                committer_scores[committer_name]['try_catch_usage']['missing_use'] += try_catch_stats['missing_use']

    for committer, skills in committer_scores.items():
        try_catch_usage = skills['try_catch_usage']
        total_usage = try_catch_usage['correctly_used'] + try_catch_usage['incorrectly_used'] + try_catch_usage['missing_use']

        if total_usage == 0:
            try_catch_usage['score'] = 'N/A'
        elif try_catch_usage['correctly_used'] == total_usage:
            try_catch_usage['score'] = 5
        elif try_catch_usage['correctly_used'] == total_usage - try_catch_usage['missing_use']:
            try_catch_usage['score'] = 4
        elif try_catch_usage['incorrectly_used'] == 0:
            try_catch_usage['score'] = 3
        elif try_catch_usage['correctly_used'] == 0:
            try_catch_usage['score'] = 2
        else:
            try_catch_usage['score'] = 1

    return committer_scores

def analyze_file(file_path):
    try_catch_stats = {
        'correctly_used': 0,
        'incorrectly_used': 0,
        'missing_use': 0
    }

    with open(file_path, 'r') as file:
        file_content = file.read()

    try:
        tree = javalang.parse.parse(file_content)
        for _, node in tree:
            if isinstance(node, javalang.tree.TryStatement):
                if node.catches:
                    try_catch_stats['correctly_used'] += 1
                else:
                    try_catch_stats['incorrectly_used'] += 1
            elif isinstance(node, javalang.tree.MethodInvocation) or isinstance(node, javalang.tree.BinaryOperation) or isinstance(node, javalang.tree.ArraySelector) or isinstance(node, javalang.tree.Cast):
                try_catch_stats['missing_use'] += 1
    except javalang.parser.JavaSyntaxError:
        pass

    return try_catch_stats

def main():
    repo_path = '/Users/damianmarquez/Documents/up/fury_ads-search-charge-bonification'
    committer_scores = analyze_try_catch_usage(repo_path)
    print(committer_scores)

if __name__ == '__main__':
    main()