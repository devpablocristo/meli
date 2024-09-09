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
    try_blocks = []
    potential_exceptions = []

    for path, node in method:
        if isinstance(node, javalang.tree.TryStatement):
            try_blocks.append(node)
        
        # Identificar potenciales situaciones de excepción
        if isinstance(node, javalang.tree.MethodInvocation):
            if node.member in ['read', 'write', 'open', 'close']:
                potential_exceptions.append(('IOException', node))
            elif node.member in ['parseInt', 'parseDouble']:
                potential_exceptions.append(('NumberFormatException', node))
        elif isinstance(node, javalang.tree.BinaryOperation) and node.operator == '/':
            potential_exceptions.append(('ArithmeticException', node))
        elif isinstance(node, javalang.tree.ArraySelector):
            potential_exceptions.append(('ArrayIndexOutOfBoundsException', node))
        elif isinstance(node, javalang.tree.Cast):
            potential_exceptions.append(('ClassCastException', node))

    # Analizar try-catch existentes
    for try_block in try_blocks:
        analyze_try_block(try_block, potential_exceptions, file_path, method_name)
    
    # Verificar excepciones potenciales no manejadas
    unhandled_exceptions = [exc for exc, node in potential_exceptions if not is_exception_handled(exc, try_blocks)]
    if unhandled_exceptions:
        print(f"Falta try-catch en {file_path}, método {method_name}:")
        for exc in set(unhandled_exceptions):
            print(f"  - Considerar manejar {exc}")

def analyze_try_block(try_block, potential_exceptions, file_path, method_name):
    caught_exceptions = [catch.parameter.types[0] for catch in try_block.catches if catch.parameter.types]
    
    # Verificar si el try-catch está bien aplicado
    relevant_exceptions = [exc for exc, node in potential_exceptions if is_node_in_try(node, try_block)]
    if relevant_exceptions:
        print(f"Try-catch bien aplicado en {file_path}, método {method_name}:")
        for exc in set(relevant_exceptions):
            print(f"  - Maneja correctamente {exc}")
    
    # Verificar si hay excepciones capturadas innecesariamente
    unnecessary_exceptions = set(caught_exceptions) - set(relevant_exceptions)
    if unnecessary_exceptions:
        print(f"Try-catch potencialmente mal aplicado en {file_path}, método {method_name}:")
        for exc in unnecessary_exceptions:
            print(f"  - {exc} capturada pero no parece necesaria")

def is_exception_handled(exception, try_blocks):
    return any(exception in [catch.parameter.types[0] for catch in block.catches if catch.parameter.types] for block in try_blocks)

def is_node_in_try(node, try_block):
    # Esta es una simplificación. En un caso real, necesitarías comparar las posiciones en el código
    return node in try_block.block

# Uso del script
repo_path = '/Users/damianmarquez/Documents/GitHub/MeliEjemplos/JAVA/fury_ads-search-test-users/.git'
analyze_try_catch(repo_path)