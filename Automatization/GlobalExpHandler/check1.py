import os
from git import Repo
import javalang
import re

def analyze_exception_handling(repo_path):
    repo = Repo(repo_path)
    
    global_handler = None
    custom_exceptions = set()
    thrown_exceptions = set()
    handled_exceptions = set()

    for blob in repo.head.commit.tree.traverse():
        if blob.path.endswith('.java'):
            file_content = blob.data_stream.read().decode('utf-8')
            
            try:
                tree = javalang.parse.parse(file_content)
                
                if is_global_exception_handler(tree):
                    global_handler = analyze_global_handler(tree, file_content)
                else:
                    analyze_file_for_exceptions(tree, custom_exceptions, thrown_exceptions)
                
            except javalang.parser.JavaSyntaxError:
                print(f"Error al parsear {blob.path}")

    if global_handler:
        handled_exceptions = set(global_handler['handled_exceptions'])
        print("Global Exception Handler encontrado:")
        print(f"Clase: {global_handler['class_name']}")
        print("Excepciones manejadas:")
        for exc in handled_exceptions:
            print(f"  - {exc}")
    else:
        print("No se encontró un Global Exception Handler")

    print("\nExcepciones personalizadas definidas:")
    for exc in custom_exceptions:
        print(f"  - {exc}")

    print("\nExcepciones lanzadas en el código:")
    for exc in thrown_exceptions:
        print(f"  - {exc}")

    unhandled_exceptions = thrown_exceptions - handled_exceptions
    if unhandled_exceptions:
        print("\nExcepciones lanzadas pero no manejadas globalmente:")
        for exc in unhandled_exceptions:
            print(f"  - {exc}")
    else:
        print("\nTodas las excepciones lanzadas están siendo manejadas globalmente.")

def is_global_exception_handler(tree):
    for _, node in tree:
        if isinstance(node, javalang.tree.ClassDeclaration):
            for decorator in node.decorators:
                if decorator.name == 'ControllerAdvice' or decorator.name == 'RestControllerAdvice':
                    return True
    return False

def analyze_global_handler(tree, file_content):
    handler_info = {'class_name': None, 'handled_exceptions': []}
    
    for _, node in tree:
        if isinstance(node, javalang.tree.ClassDeclaration):
            handler_info['class_name'] = node.name
            break
    
    # Usar expresiones regulares para encontrar los métodos anotados con @ExceptionHandler
    exception_handler_pattern = r'@ExceptionHandler\s*\(\s*(\w+)\.class\s*\)'
    matches = re.findall(exception_handler_pattern, file_content)
    handler_info['handled_exceptions'] = matches

    return handler_info

def analyze_file_for_exceptions(tree, custom_exceptions, thrown_exceptions):
    for path, node in tree:
        if isinstance(node, javalang.tree.ClassDeclaration):
            if 'Exception' in node.extends:
                custom_exceptions.add(node.name)
        elif isinstance(node, javalang.tree.MethodDeclaration):
            for _, child_node in node:
                if isinstance(child_node, javalang.tree.ThrowStatement):
                    if child_node.throws:
                        thrown_exceptions.add(child_node.throws.member)
        elif isinstance(node, javalang.tree.TryStatement):
            for catch in node.catches:
                thrown_exceptions.add(catch.parameter.types[0])

# Uso del script
repo_path = '/Users/damianmarquez/Documents/GitHub/MeliEjemplos/JAVA/fury_ads-search-test-users/.git'
analyze_exception_handling(repo_path)