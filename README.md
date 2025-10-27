# 🗃️ Projeto Pandora

> “Quando abrimos a caixa, descobrimos... que o Go não era tão rápido assim.” 🐢⚡🦀

## 📖 Sobre o Projeto

O **Pandora** nasceu como uma **prova de conceito** escrita em **Go (Golang)** com o objetivo de testar desempenho, organização e escalabilidade em um serviço de API moderna.

A proposta inicial era simples: verificar o quanto o Go aguentaria sob carga intensa de I/O e consultas complexas — especialmente comparado ao **Rust**, que sempre se vendeu como o queridinho da performance e segurança de memória.

## ⚙️ Stack Original (Go)

- **Linguagem:** Go 1.25+
- **Framework HTTP:** [Fiber](https://gofiber.io/)
- **ORM:** GORM
- **Banco:** PostgreSQL
- **Cache:** Redis
- **Arquitetura:** Controller → Service → Repository
- **Dockerizado:** Sim (imagem baseada em `alpine`)

## 🔬 Resultado dos Testes

Após algumas rodadas de benchmark e profiling… digamos que o **Rust** deu um _“cacete épico”_ no Go em praticamente todos os cenários testados 😅:

| Cenário                        | Go     | Rust   | Diferença                           |
| ------------------------------ | ------ | ------ | ----------------------------------- |
| Requisições simultâneas (100k) | 🐢     | ⚡     | Rust ≈ 3x mais rápido               |
| Consumo de memória             | 📈     | 📉     | Rust consumiu ~40% menos            |
| Latência média                 | ~4.8ms | ~1.6ms | Rust manteve estabilidade sob carga |

O veredito foi claro: **o Go ficou para trás nessa corrida específica.**

## 🔁 Migração para Rust

Diante dos resultados, o projeto foi **totalmente migrado para Rust**, utilizando:

- **Framework:** Actix Web
- **ORM:** SQLx
- **Cache:** Redis
- **Banco:** PostgreSQL
- **Build:** Cargo

A reescrita mostrou ganhos consideráveis em:

- Tempo de resposta
- Controle fino de memória
- Melhor paralelismo via `tokio`

## 💬 Conclusão

O **Go** continua sendo uma linguagem incrível — rápida de desenvolver (e como é rapidá, é incrivel oque se pode fazer em algumas horas dedicadas com go), produtiva e com uma comunidade sólida.  
Mas para este tipo de aplicação, altamente dependente de performance e controle de baixo nível, o **Rust simplesmente dominou**.

> “Nem sempre a caixa de Pandora guarda apenas males — às vezes, ela revela um novo amor pela performance.” ❤️🦀

---
