# Remove Background Desktop App

Aplicação desktop para remoção de fundo de imagens, 100% offline, construída com Wails v2, Go e Vue 3.

## Arquitetura Geral

```
┌─────────────────────────────────────────────────────────────┐
│                     Desktop App                             │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────────┐      ┌─────────────────────────────┐ │
│  │   Vue 3 Frontend │◄────►│      Go Backend             │ │
│  │   (Renderer)    │ IPC  │    (Wails Bridge)           │ │
│  └─────────────────┘      └─────────────────────────────┘ │
│                                      │                     │
│                               ┌──────▼──────┐             │
│                               │ ONNX Runtime│             │
│                               │   (CGO)     │             │
│                               └─────────────┘             │
└─────────────────────────────────────────────────────────────┘
```

### Stack Tecnológico

| Camada            | Tecnologia             |
| ----------------- | ---------------------- |
| Desktop Framework | Wails v2               |
| Backend           | Go 1.23+               |
| Frontend          | Vue 3 + TypeScript     |
| UI                | TailwindCSS v3         |
| State Management  | Pinia                  |
| Inferência ML     | ONNX Runtime (via CGO) |
| Modelo            | RMBG-2.0 (BiRefNet)    |

---

## Estrutura de Diretórios

```
remove-bg-go/
├── app.go                    # Aplicação principal Wails
├── main.go                   # Entry point
├── go.mod                    # Dependências Go
├── wails.json                # Configuração Wails
├── Makefile                  # Comandos de build
│
├── internal/
│   └── inference/
│       └── inference.go       # Módulo de inferência ONNX
│
├── models/
│   └── RMBG-2.0/
│       └── onnx/
│           └── model.onnx    # Modelo de ML (BiRefNet)
│
├── frontend/
│   ├── src/
│   │   ├── App.vue                    # Componente raiz
│   │   ├── main.ts                    # Entry point Vue
│   │   ├── style.css                  # Estilos globais
│   │   ├── components/
│   │   │   ├── TitleBar.vue           # Titlebar customizada
│   │   │   ├── DropZone.vue           # Área de upload
│   │   │   ├── ProcessingView.vue     # Loading/Progresso
│   │   │   ├── ResultView.vue         # Resultado com slider
│   │   │   └── Toast.vue              # Notificações
│   │   ├── stores/
│   │   │   └── app.ts                 # Estado global (Pinia)
│   │   └── lib/
│   │       ├── design-tokens.ts       # Tokens de design
│   │       └── utils.ts               # Utilitários
│   │
│   ├── wailsjs/                       # Bindings Go→JS (gerado)
│   │   ├── go/main/App.js             # Wrapper para métodos Go
│   │   └── runtime/runtime.js         # Runtime Wails
│   │
│   ├── package.json
│   ├── vite.config.ts
│   └── tailwind.config.js
│
└── build/                    # Saída do build
```

---

## Fluxo de Dados

### 1. Seleção de Imagem

```
Usuário arrasta/seleciona imagem
         │
         ▼
DropZone.vue: handleFile()
         │
         ├─ Valida tipo (PNG/JPG/WEBP)
         ├─ Valida tamanho (<20MB)
         ├─ Converte para Base64
         │
         ▼
store.setOriginalImage(base64)
store.setLoading()
         │
         ▼
App.vue: Transição para ProcessingView
```

### 2. Processamento no Backend

```
DropZone.vue: RemoveBackground(imageBase64)
         │
         ▼
Wails IPC Bridge
         │
         ▼
app.go: App.RemoveBackground()
         │
         ├─ 1. Carrega modelo ONNX (se primeira vez)
         │      └─ inference.NewSession()
         │
         ├─ 2. Decodifica imagem Base64 → image.Image
         │
         ├─ 3. Pré-processamento
         │      └─ inference.Preprocessor.Preprocess()
         │            ├─ Resize para 1024x1024
         │            └─ Normalização: mean=[0.5,0.5,0.5], std=[0.5,0.5,0.5]
         │
         ├─ 4. Inferência ONNX
         │      └─ session.RunInference()
         │            Input:  [1, 3, 1024, 1024] float32
         │            Output: [1, 1, 1024, 1024] float32
         │
         ├─ 5. Pós-processamento
         │      └─ inference.Postprocessor.Postprocess()
         │            ├─ Redimensiona máscara para tamanho original
         │            ├─ Aplica alpha channel
         │            └─ Codifica para PNG
         │
         ▼
Retorna Base64 PNG com transparência
```

### 3. Exibição do Resultado

```
RemoveBackground() retorna resultBase64
         │
         ▼
store.setResultImage(resultBase64)
store.setDone()
         │
         ▼
App.vue: Transição para ResultView
         │
         ├─ Exibe comparação: original | resultado
         ├─ Slider para comparar antes/depois
         ├─ Botão Download PNG
         └─ Botão "Remove another"
```

---

## Componentes Frontend

### TitleBar.vue

- Titlebar customizada (frameless window)
- Região de drag (`--wails-draggable: drag`)
- Botões de controle de janela (minimize, maximize, close)

### DropZone.vue

- Área de drop com validação
- Animação de hover com glow
- Converte arquivo para Base64
- **Chama o método Go RemoveBackground**

### ProcessingView.vue

- Exibe imagem original
- Barra de progresso animada
- Mensagens de status via EventsOn

### ResultView.vue

- Layout split (original | resultado)
- **Slider de comparação** customizado (drag-to-reveal)
- Background checkerboard para mostrar transparência
- Botões de ação: Download, Copy, Reset

### Toast.vue

- Notificações de erro/sucesso
- Auto-dismiss após 4 segundos

---

## Estado Global (Pinia)

Arquivo: [`frontend/src/stores/app.ts`](frontend/src/stores/app.ts)

```typescript
type AppState = "idle" | "loading" | "processing" | "done" | "error";

interface State {
  currentState: AppState;
  originalImage: string | null; // Base64 da imagem original
  resultImage: string | null; // Base64 do resultado
  statusMessage: string; // Mensagem de status do Go
  error: string | null; // Mensagem de erro
  history: HistoryItem[]; // Últimas 5 imagens processadas
  progress: number; // Progresso (0-100)
}
```

---

## Backend Go

### app.go

Ponto de entrada principal do Wails. Expõe métodos para o frontend:

```go
type App struct {
    ctx           context.Context
    session       *inference.Session      // ONNX session (singleton)
    preprocessor  *inference.Preprocessor
    postprocessor *inference.Postprocessor
    modelPath     string
}

// Métodos expostos ao frontend
func (a *App) RemoveBackground(imageBase64 string) (string, error)
func (a *App) GetVersion() string
```

### internal/inference/inference.go

Módulo de inferência ONNX:

| Tipo            | Descrição                                            |
| --------------- | ---------------------------------------------------- |
| `Session`       | Gerencia sessão ONNX Runtime com tensor pré-alocados |
| `Preprocessor`  | Redimensiona e normaliza imagens                     |
| `Postprocessor` | Converte máscara em imagem com alpha                 |

#### Session (Singleton Thread-Safe)

```go
type Session struct {
    session     *ort.AdvancedSession    // Sessão ONNX
    inputTensor  *ort.Tensor[float32]   // Tensor de entrada pré-alocado
    outputTensor *ort.Tensor[float32]   // Tensor de saída pré-alocado
}
```

- **Carrega o modelo uma única vez** no boot
- **Reutiliza tensores** entre chamadas (sem overhead de alocação)
- Tenta múltiplos nomes de input/output comuns:
  - `input` / `sigmoid_0`
  - `x` / `sigmoid`
  - `input.1` / `output.1`

#### Preprocess

1. Resize para 1024x1024 (nearest neighbor)
2. Normalização por canal:
   - `normalized = (pixel / 255.0 - mean) / std`
   - `mean = [0.5, 0.5, 0.5]`
   - `std = [0.5, 0.5, 0.5]`
3. Converte para formato NCHW: `[1, 3, 1024, 1024]`

#### Postprocess

1. Normaliza saída para 0-255
2. Resize da máscara para dimensões originais
3. Aplica alpha channel na imagem original
4. Codifica para PNG com transparência

---

## Comunicação Frontend ↔ Backend

### Chamada de Método

O Wails gera automaticamente bindings em [`frontend/wailsjs/go/main/App.js`](frontend/wailsjs/go/main/App.js):

```javascript
// Gerado automaticamente pelo Wails
export function RemoveBackground(arg1) {
  return window["go"]["main"]["App"]["RemoveBackground"](arg1);
}
```

Uso no Vue:

```typescript
import { RemoveBackground } from "@wailsjs/go/main/App";

const result = await RemoveBackground(base64Image);
```

### Eventos (Server → Client)

O Go emite eventos de status via `runtime.EventsEmit`:

```go
runtime.EventsEmit(a.ctx, "status", "loading_model")
runtime.EventsEmit(a.ctx, "status", "decoding")
runtime.EventsEmit(a.ctx, "status", "preprocessing")
runtime.EventsEmit(a.ctx, "status", "processing")
runtime.EventsEmit(a.ctx, "status", "finalizing")
runtime.EventsEmit(a.ctx, "status", "done")
```

O frontend ouve com `EventsOn`:

```typescript
import { EventsOn } from "@wailsjs/runtime";

EventsOn("status", (status: string) => {
  console.log("Status:", status);
});
```

---

## Design Tokens

Arquivo: [`frontend/src/lib/design-tokens.ts`](frontend/src/lib/design-tokens.ts)

```typescript
export const colors = {
  background: "#0A0A0F", // Preto com tom frio
  surface: "#111118", // Cards
  surfaceElevated: "#1C1C27", // Hover states
  border: "#2A2A3D",
  accent: "#7C6EFA", // Violeta/Índigo
  accentHover: "#9D8FFF",
  success: "#34D399",
  error: "#F87171",
  textPrimary: "#F0F0FF",
  textMuted: "#8B8BA7",
};
```

---

## Build e Execução

### Development

```bash
# Instalar dependências Go
go mod tidy

# Rodar em modo desenvolvimento
wails dev
```

### Production Build

```bash
# Build completo
wails build

# Ou usando Makefile
make build
```

O binário será gerado em: `build/bin/remove-bg-go`

### Variáveis de Ambiente

| Variável          | Descrição                         |
| ----------------- | --------------------------------- |
| `CGO_ENABLED=1`   | Requerido para ONNX Runtime (CGO) |
| `LD_LIBRARY_PATH` | Path para libonnxruntime.so       |

---

## Dependências

### Go

- `github.com/wailsapp/wails/v2` - Framework desktop
- `github.com/yalue/onnxruntime_go` - Wrapper ONNX Runtime

### Frontend (npm)

- `vue` - Framework UI
- `pinia` - State management
- `tailwindcss` - Estilização
- `lucide-vue-next` - Ícones

---

## Modelo de ML

O app usa o modelo **RMBG-2.0** (BiRefNet) em formato ONNX.

Local: `models/RMBG-2.0/onnx/model.onnx`

O modelo:

- Input: `[1, 3, 1024, 1024]` (RGB normalizado)
- Output: `[1, 1, 1024, 1024]` (máscara de probabilidade)
- Taille: ~170MB (versão FP16)

---

## Limitações e Considerações

1. **CGO Required**: O ONNX Runtime requer CGO, então o build é vinculado à plataforma
2. **Memória**: O modelo consome ~500MB de RAM durante inferência
3. **Tamanho de imagem**: Otimizado para imagens de até 20MB
4. **Formatos**: Suporta apenas PNG, JPG, WEBP
