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
                
                # Analizar el árbol sintáctico
                analyze_tree(tree, blob.path)
                
            except javalang.parser.JavaSyntaxError:
                print(f"Error al parsear {blob.path}")

def analyze_tree(tree, file_path):
    try_blocks = []
    method_declarations = []
    
    for path, node in tree:
        if isinstance(node, javalang.tree.TryStatement):
            try_blocks.append(node)
        elif isinstance(node, javalang.tree.MethodDeclaration):
            method_declarations.append(node)
    
    # Verificar métodos sin try-catch
    for method in method_declarations:
        if not any(try_block in method.children for try_block in try_blocks):
            print(f"Advertencia: Método sin try-catch en {file_path}: {method.name}")
    
    # Analizar los bloques try-catch
    for try_block in try_blocks:
        if not try_block.catches:
            print(f"Error: Bloque try sin catch en {file_path}")
        elif any(catch.parameter.types == [] for catch in try_block.catches):
            print(f"Advertencia: Catch genérico (Exception) en {file_path}")

# Uso del script
repo_path = '/Users/damianmarquez/Documents/GitHub/MeliEjemplos/JAVA/fury_ads-search-test-users/.git'
analyze_try_catch(repo_path)