import os
from git import Repo
import javalang

# Definir las capas esperadas en la arquitectura hexagonal
EXPECTED_LAYERS = {
    'domain': ['model', 'service', 'port'],
    'application': ['usecase'],
    'infrastructure': ['adapter', 'config', 'repository']
}

def analyze_hexagonal_architecture(repo_path):
    repo = Repo(repo_path)
    
    # Estructura para almacenar los resultados del análisis
    architecture_structure = {layer: {sublayer: [] for sublayer in sublayers} for layer, sublayers in EXPECTED_LAYERS.items()}
    
    for blob in repo.head.commit.tree.traverse():
        if blob.path.endswith('.java'):
            file_path = blob.path
            file_content = blob.data_stream.read().decode('utf-8')
            
            try:
                tree = javalang.parse.parse(file_content)
                analyze_file(tree, file_path, architecture_structure)
            except javalang.parser.JavaSyntaxError:
                print(f"Error al parsear {file_path}")
    
    validate_architecture(architecture_structure)

def analyze_file(tree, file_path, architecture_structure):
    package_name = next((decl.name for decl in tree.package.imports if isinstance(decl, javalang.tree.PackageDeclaration)), None)
    
    if not package_name:
        return
    
    for layer, sublayers in EXPECTED_LAYERS.items():
        if layer in package_name:
            for sublayer in sublayers:
                if sublayer in package_name:
                    for _, node in tree:
                        if isinstance(node, javalang.tree.ClassDeclaration):
                            architecture_structure[layer][sublayer].append((node.name, file_path))
                    break

def validate_architecture(architecture_structure):
    print("Análisis de Arquitectura Hexagonal:")
    
    # Verificar la presencia de las capas principales
    for layer in EXPECTED_LAYERS:
        if any(architecture_structure[layer].values()):
            print(f"✓ Capa '{layer}' presente")
        else:
            print(f"✗ Capa '{layer}' ausente")
    
    # Verificar la estructura dentro de cada capa
    for layer, sublayers in EXPECTED_LAYERS.items():
        print(f"\nAnálisis de la capa '{layer}':")
        for sublayer in sublayers:
            classes = architecture_structure[layer][sublayer]
            if classes:
                print(f"  ✓ Sublayer '{sublayer}' contiene {len(classes)} clases")
                for class_name, file_path in classes:
                    print(f"    - {class_name} ({file_path})")
            else:
                print(f"  ✗ Sublayer '{sublayer}' vacía")
    
    # Verificar dependencias (simplificado)
    print("\nVerificación de dependencias (simplificado):")
    if architecture_structure['domain']['port']:
        print("✓ Puertos definidos en el dominio")
    else:
        print("✗ No se encontraron puertos en el dominio")
    
    if architecture_structure['infrastructure']['adapter']:
        print("✓ Adaptadores definidos en la infraestructura")
    else:
        print("✗ No se encontraron adaptadores en la infraestructura")

# Uso del script
repo_path = '/Users/damianmarquez/Documents/GitHub/MeliEjemplos/JAVA/fury_ads-search-test-users/.git'
analyze_hexagonal_architecture(repo_path)