Read in [English](README.md)

# Smart Transaction Repair API

![technology Go](https://img.shields.io/badge/technology-go-blue.svg)

## Objetivo:

É permitir ao usuário, de forma simples e direta, solicitar a devolução do dinheiro de uma transação que foi negada, porém, cobrada.

> ### Sites Suportados:
>
> | Site | 
> |------------|
> | México (MLM) | 

---

</br>

## Devolução:

O pedido de devolução passará por uma validação, com regras previamente definidas, que validará a elegibilidade do usuário solicitante.


> ### Regras Consideradas:
>
> | Parâmetro | Regra 
> |----------|-----|
> | max_amount_reparation | o valor da transação não deve ser maior que o valor máximo permitido |
> | qty_reparation_per_period_days | há um limite para qtde de devoluções dentro de um intervalo de dias |
> | status_detail_allowed | o status da transação deve estar dentro do que é permitido |
> | usuário | o usuário não deve estar bloqueado** | 
>

>
> ### Valores Parametrizáveis 
> 
> Os valores definidos que serão utilizados para aplicação das regras, estão armazenados em formato `json` dentro do serviço de configuração do Fury `Configurations`, permitindo a alteração dos valores sem a necessidade de alteração no código.
>
> Exemplo json:
> ``` json
> { 
>   "status_detail_allowed": {
>     "pending_capture": {}
>   },
>   "qty_reparation_per_period_days": {
>     "qty": 2,
>     "period_days": 30
>   },  
>   "max_amount_reparation": 200000
> }
>
> ```
> ** A verificação de um usuário bloqueado, é feita diretamente na base de dados blockedlist.

---

</br>

## Bloqueio de Usuário:

 O usuário é bloqueado quando o mesmo realiza uma devolução e, em seguida a captura da transação devolvida é transmitida para o nosso eco sistema.

 A verificação é feita através de um consumidor do tópico `cards-transactions-api-news-feed.cards-transactions-api`, onde para cada captura recebida, a aplicação verifica se existe devolução associada à captura através do identificador da autorização. 

---

</br>

## API

### Ambientes:
| Ambiente | Protocolo | Host | BasePath
|----------|-----|-|-|
| beta | http | internal.mercadopago.com | /beta/cards/smart-transaction-repair |
| prod | http | internal.mercadopago.com | /cards/smart-transaction-repair |

</br>

### Rotas:

### Solicitar Devolução de uma Transação

> #### POST /v1/reverse/{payment_id}

##### Request Path Parameters

| Parâmetro | Descrição |
|-----------|-----------|
| payment_id | identificador do pagamento |

#####  Header

| Parâmetro | Descrição |
|-----------|-----------|
| X-Client-Id | identificador da aplicação solicitante |

</br>

##### Request Body

``` json
{
    "user_id":123
}

```

| Propriedade | Descrição | Tipo | Obrigatório |
|-------------|-----------|------|------------|
| user_id | identificador do usuário solicitante da devolução | int | SIM |

</br>

##### 1. Response Body

``` json
{
    "message":"Reverse successfully requested"
}

```

| Propriedade | Descrição | Tipo |
|-------------|-----------|------|
| message | mensagem de sucesso da devolução | string |

</br>

##### 1.1. Response Body Error

```json
{
    "code": "unauthorized",
    "message": "invalid request",
    "cause": "request is not authorized"
}
```

| Propriedade | Descrição | Tipo |
|-------------|-----------|------|
| code | código de erro da aplicação | string |
| message | mensagem genérica de erro | string |
| cause | motivo do erro | string |

</br>

##### 1.2. Response Body Error - Not Eligible

```json
{
    "code": "not_eligible",
    "message": "validation result",
    "cause": {
        "reason": "customer not eligible for reversal",
        "creation_datetime": "2022-12-07T21:46:07.713Z"         
    }
}
```

| Propriedade | Descrição | Tipo |
|-------------|-----------|------|
| code | código de erro da aplicação | string |
| message | mensagem genérica de erro | string |
| cause | detalhes do erro | object |
| cause.reason | mensagem genérica da validação da devolução negada | string |
| cause.creation_datetime | data de criação do pagamento | string |

</br>

---

## Arquitetura
![architecture](resources/architecture.png)

</br>

---
## Diagrama Sequência > [fluxos](workflow.pt.md?id=diagrama-sequencia)

</br>

---
## Queries > [exemplos](queries.md?id=queries)