## GitFlow

1. **`$ gco develop`**
   - **Comando original**: `git checkout develop`
   - **Descrição**: Muda para a branch chamada "develop", uma prática comum no GitFlow para trabalhar com a linha principal do desenvolvimento.

2. **`$ ggl`**
   - **Comando original**: `git pull`
   - **Descrição**: Atualiza a branch local com as mudanças mais recentes do repositório remoto, garantindo que você está trabalhando com a versão mais atual.

3. **`$ gcb feature/nova-branch`**
   - **Comando original**: `git checkout -b feature/nova-branch`
   - **Descrição**: Cria e muda para uma nova branch chamada "feature/nova-branch", utilizada para desenvolver novas funcionalidades de forma isolada.

4. **`$ gaa`**
   - **Comando original**: `git add --all`
   - **Descrição**: Adiciona todos os arquivos modificados à área de preparação (staging area) para o próximo commit, um passo fundamental antes de efetuar um commit.

5. **`$ gcmsg "novo commit"`**
   - **Comando original**: `git commit -m "novo commit"`
   - **Descrição**: Cria um novo commit com a mensagem "novo commit", registrando oficialmente as mudanças no histórico do repositório.

6. **`$ gf`**
   - **Comando original**: `git fetch`
   - **Descrição**: Busca as últimas mudanças e atualizações das branches no repositório remoto, sem mesclar com a branch local.

7. **`$ gm develop`**
   - **Comando original**: `git merge develop`
   - **Descrição**: Funde a branch "develop" na branch atual. Se houver conflitos, é necessário resolvê-los antes de concluir a fusão.

8. **`$ ggp`**
   - **Comando original**: `git push`
   - **Descrição**: Envia as alterações da branch local para o repositório remoto, compartilhando seus commits com outros colaboradores.

9. **Abra o repositório no GitHub**: Visualize as branches e as mudanças na interface web do GitHub.

10. **Criar um Pull Request**: Use o aviso para criar um Pull Request para sua nova branch; clique em "Comparar & solicitar pull".

11. **Página de "Abrir um pull request"**: Escreva um título e uma descrição para o seu PR, selecionando a branch para a qual deseja que sua branch seja fundida (geralmente é a branch main ou develop).

12. **"Criar pull request"**: Inicie o processo de revisão pelo time.

13. **Revisão do PR**: Outros membros da equipe podem revisar e comentar.

14. **Atualizações**: Faça mudanças adicionais se necessário; o PR atualiza-se automaticamente.

15. **Aprovação do PR**: Uma vez aprovado, alguém com permissões pode mesclar sua branch na branch destino.

### Resolvendo Conflitos no Git

1. **Identificar os Conflitos**: Use `git status` para identificar arquivos não mesclados (conflitantes).

2. **Revisar os Conflitos**: Abra os arquivos indicados e procure por marcas de conflito (<<<<<<<, =======, >>>>>>>).

3. **Editar os Arquivos para Resolver Conflitos**: Decida o conteúdo correto para cada seção conflitante e edite conforme necessário.

4. **Marcar os Arquivos como Resolvidos**: Após resolver os conflitos em um arquivo, use `git add [nome-do-arquivo]` para marcá-lo como resolvido.

5. **Finalizar a Fusão ou Rebase**: Complete o processo com `git commit` para um merge.

6. **Verificar e Testar**: Assegure-se de que o código funciona como esperado após a resolução dos conflitos.

7. **Subir Mudanças ao Repositório Remoto**: Use `git push` para enviar as alterações para o remoto.

### Notas Adicionais

- **`git pull origin main`**: Uma forma eficaz de sincronizar o branch de trabalho local com o branch "main" no repositório remoto.

- **`git`**: O comando base para o sistema de controle de versões Git.

- **`pull`**: Um comando Git para atualizar

 o repositório local com mudanças do remoto, combinando `git fetch` e `git merge`.

- **`origin`**: O nome padrão para o repositório remoto de onde o local foi clonado.

- **`main`**: O branch principal em muitos repositórios, contendo o código base mais atualizado.