# [Pós GoExpert - FullCycle](https://fullcycle.com.br)

## Desafios técnicos - Deploy com Cloud Run

### Pré-requisitos
- [Golang](https://golang.org/)

### Como executar a aplicação

```bash
  # 1 - Clonar o repositório do projeto
  git clone https://github.com/andersonigorf/goexpert-cloud-run.git
  
  # 2 - Acessar o diretório do projeto
  cd goexpert-cloud-run

  # 3 - Executar a aplicação com o Docker
  docker-compose up -d
  
  ou
  
  make run
  
  # 4 - Executar testes automatizados
  go test ./... -v
  
  ou
  
  make test
  
```

### Exemplos de requisições

```bash
  # 1 - Executar as requisições do arquivo requests.http (dentro da pasta ./api)
  api/requests.http
  
  # 2 - Executar pelo browser
  http://localhost:8080/?cep=<CEP>
   
  http://localhost:8080/?cep=72547240

  # 3 - Executar pelo Endereço do serviço no Google Cloud Run
  https://cloudrun-goexpert-z35ax4go6q-uc.a.run.app/?cep=<CEP>
  
  https://cloudrun-goexpert-z35ax4go6q-uc.a.run.app/?cep=72547240
```