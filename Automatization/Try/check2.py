import os
from git import Repo
import javalang

def analyze_try_catch(repo_path):
    repo = Repo(repo_path)
    
    for blob in repo.head.commit.tree.traverse():
        if blob.path.endswith('.java'):
            file_content = blob.data_stream.read().decode('utf-8')
            
            try:
                tree = javalang.parse.parse(file_content)
                analyze_tree(tree, blob.path)
            except javalang.parser.JavaSyntaxError:
                print(f"Error al parsear {blob.path}")

def analyze_tree(tree, file_path):
    for path, node in tree:
        if isinstance(node, javalang.tree.MethodDeclaration):
            analyze_method(node, file_path)

def analyze_method(method, file_path):
    method_name = method.name
    has_try_catch = False
    potential_exceptions = []

    for path, node in method:
        if isinstance(node, javalang.tree.TryStatement):
            has_try_catch = True
        
        # Verificar llamadas a métodos que podrían lanzar excepciones
        if isinstance(node, javalang.tree.MethodInvocation):
            if node.member in ['read', 'write', 'open', 'close']:
                potential_exceptions.append('IOException')
            elif node.member in ['parseInt', 'parseDouble']:
                potential_exceptions.append('NumberFormatException')
        
        # Verificar operaciones de división
        if isinstance(node, javalang.tree.BinaryOperation) and node.operator == '/':
            potential_exceptions.append('ArithmeticException')
        
        # Verificar accesos a arrays
        if isinstance(node, javalang.tree.ArraySelector):
            potential_exceptions.append('ArrayIndexOutOfBoundsException')
        
        # Verificar casteos
        if isinstance(node, javalang.tree.Cast):
            potential_exceptions.append('ClassCastException')

    if potential_exceptions and not has_try_catch:
        print(f"Sugerencia: Considerar agregar try-catch en el método {method_name} en {file_path}")
        print(f"  Posibles excepciones: {', '.join(set(potential_exceptions))}")

# Uso del script
repo_path = '/Users/damianmarquez/Documents/GitHub/MeliEjemplos/JAVA/fury_ads-search-test-users/.git'
analyze_try_catch(repo_path)