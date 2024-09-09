import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.HashMap;
import java.util.Map;
import java.util.stream.Collectors;

public class HexagonalArchitectureAnalyzer {
    private static final String[] EXPECTED_LAYERS = {"domain", "application", "infrastructure"};
    private static final String[] EXPECTED_SUBLAYERS = {"model", "service", "port", "usecase", "adapter", "config", "repository"};

    public static void main(String[] args) {
        if (args.length < 1) {
            System.out.println("Por favor, proporciona la ruta al repositorio como argumento.");
            return;
        }

        String repoPath = "/Users/damianmarquez/Documents/GitHub/MeliEjemplos/JAVA/fury_ads-search-test-users/.git";
        analyzeHexagonalArchitecture(repoPath);
    }

    private static void analyzeHexagonalArchitecture(String repoPath) {
        Map<String, Map<String, java.util.List<String>>> architectureStructure = new HashMap<>();
        for (String layer : EXPECTED_LAYERS) {
            Map<String, java.util.List<String>> sublayerStructure = new HashMap<>();
            for (String sublayer : EXPECTED_SUBLAYERS) {
                sublayerStructure.put(sublayer, new java.util.ArrayList<>());
            }
            architectureStructure.put(layer, sublayerStructure);
        }

        try {
            Files.walk(Paths.get(repoPath))
                 .filter(path -> path.toString().endsWith(".java"))
                 .forEach(path -> analyzeFile(path, architectureStructure));

            validateArchitecture(architectureStructure);
        } catch (IOException e) {
            System.out.println("Error al analizar el repositorio: " + e.getMessage());
        }
    }

    private static void analyzeFile(Path filePath, Map<String, Map<String, java.util.List<String>>> architectureStructure) {
        String packageName = getPackageName(filePath, Paths.get(filePath.toString().replaceAll("\\.java$", "")));

        for (String layer : EXPECTED_LAYERS) {
            if (packageName.contains(layer)) {
                for (String sublayer : EXPECTED_SUBLAYERS) {
                    if (packageName.contains(sublayer)) {
                        architectureStructure.get(layer).get(sublayer).add(filePath.getFileName().toString());
                        break;
                    }
                }
            }
        }
    }

    private static void validateArchitecture(Map<String, Map<String, java.util.List<String>>> architectureStructure) {
        System.out.println("Análisis de Arquitectura Hexagonal:");

        // Verificar la presencia de las capas principales
        for (String layer : EXPECTED_LAYERS) {
            if (architectureStructure.get(layer).values().stream().flatMap(java.util.List::stream).collect(Collectors.toList()).size() > 0) {
                System.out.println("✓ Capa '" + layer + "' presente");
            } else {
                System.out.println("✗ Capa '" + layer + "' ausente");
            }
        }

        // Verificar la estructura dentro de cada capa
        for (String layer : EXPECTED_LAYERS) {
            System.out.println("\nAnálisis de la capa '" + layer + "':");
            for (String sublayer : EXPECTED_SUBLAYERS) {
                java.util.List<String> classes = architectureStructure.get(layer).get(sublayer);
                if (classes.size() > 0) {
                    System.out.println("  ✓ Sublayer '" + sublayer + "' contiene " + classes.size() + " clases");
                    for (String className : classes) {
                        System.out.println("    - " + className);
                    }
                } else {
                    System.out.println("  ✗ Sublayer '" + sublayer + "' vacía");
                }
            }
        }

        // Verificar dependencias (simplificado)
        System.out.println("\nVerificación de dependencias (simplificado):");
        if (architectureStructure.get("domain").get("port").size() > 0) {
            System.out.println("✓ Puertos definidos en el dominio");
        } else {
            System.out.println("✗ No se encontraron puertos en el dominio");
        }

        if (architectureStructure.get("infrastructure").get("adapter").size() > 0) {
            System.out.println("✓ Adaptadores definidos en la infraestructura");
        } else {
            System.out.println("✗ No se encontraron adaptadores en la infraestructura");
        }
    }

    private static String getPackageName(Path filePath, Path baseDir) {
        return baseDir.relativize(filePath).toString().replace("/", ".");
    }
}