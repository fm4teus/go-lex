## Analisador Léxico

## Como Executar

Esse trabalho foi desenvolvido em [Go](https://go.dev/).

- Com a linguagem instalada é possível executar usando:
```
go run main.go <CAMINHO-ARQUIVO>
```

- Também é possível compilar o código com:
```
go build
```

- E executar passando o caminho para um arquivo:
```
./go-lex <CAMINHO-ARQUIVO>
```

## Estratégia

Esse analisador léxico consiste em ler um array de bytes contendo o código fonte e deve ter como saída uma tabela de tokens e também retornando erros para aqueles tokens que não puderam ser identificados.

Temos um laço for que a cada iteração pula os espaços em branco contando apenas sua posição e tenta encontrar o próximo token válido.

A função Lex cria uma thread para cada tipo de token e aguarda o retorno no canal de resposta, no caso de sucesso uma dessas threads irá retornar as informações do token encontrado, caso esse canal não receba nenhuma resposta é retornado um token de erro nessa localização.
```go
wg.Add(5)
	go matchRegexp(stringRegex, b, STRING, ch, wg)
	go matchRegexp(numberRegex, b, NUM, ch, wg)
	go matchRegexp(separatorRegex, b, SEP, ch, wg)
	go matchRegexp(identifierRegex, b, IDENTIFIER, ch, wg)
	go matchRegexp(operatorRegex, b, OPERATOR, ch, wg)
	go func() {
		wg.Wait()
		close(ch)
	}()

	res, ok := <-ch
	if !ok {
		return 1, ERROR, "Error"
	}
	return res.End, res.Token, res.Lit
``` 
Exemplo de saída:
```
:1     KEYWORD         for
2:5     NUM             4
2:7     SEPARATOR       (
2:8     KEYWORD         int
2:12    IDENTIFIER      a
2:13    OPERATOR        +=
2:16    NUM             2
2:17    SEPARATOR       ;
2:19    IDENTIFIER      a
2:20    OPERATOR        <
2:21    NUM             5
2:22    SEPARATOR       ;
2:24    IDENTIFIER      a
2:26    OPERATOR        ++
2:28    SEPARATOR       )
3:1     KEYWORD         int
3:5     IDENTIFIER      a
3:7     OPERATOR        =
3:9     NUM             2
3:11    OPERATOR        +
3:13    NUM             2
5:1     KEYWORD         if
5:4     SEPARATOR       (
5:5     IDENTIFIER      a
5:6     OPERATOR        >
5:7     NUM             2
5:9     OPERATOR        ||
5:12    IDENTIFIER      a
5:14    OPERATOR        <
5:16    OPERATOR        -
5:17    NUM             1
5:19    OPERATOR        &&
5:22    OPERATOR        -
5:23    NUM             4
```
