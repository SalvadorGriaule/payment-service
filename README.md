 Voici un README pour votre projet de service de paiement :

---

# Payment Service

Service de paiement HTTP en Go utilisant le framework Gin. Ce service gère la création et la récupération de transactions de paiement avec stockage en mémoire.

## Structure du projet

```
payment-service/
├── main.go           # Point d'entrée - configuration des routes
├── memory.go         # Package store - modèles et stockage en mémoire
└── request.go        # Package http - handlers HTTP
```

## Packages

| Package | Fichier | Description |
|---------|---------|-------------|
| `main` | `main.go` | Définition des routes et démarrage du serveur |
| `store` | `memory.go` | Structures de données et stockage en mémoire |
| `http` | `request.go` | Handlers pour les endpoints REST |

## API Endpoints

### Créer un paiement
```http
POST /v1/payments
Headers:
  X-Tenant-Id: <tenant-id>
  Idempotency-Key: <cle-idempotence>
Body:
{
    "orderRef": "ORD-123",
    "amount": 150.00,
    "currency": "EUR"
}
```

**Réponses :**
- `202 Accepted` - Paiement créé ou déjà existant (idempotence)
- `400 Bad Request` - Données invalides

**Logique métier :**
- `amount <= 0` → Status `FAILED`
- `amount >= 10000` → Status `REQUIRES_ACTION` (requiert validation supplémentaire)
- `amount < 10000` → Status `SUCCEEDED`

### Récupérer un paiement
```http
GET /v1/payments/:id
```

**Réponse :**
```json
{
    "status": "SUCCEEDED"
}
```

## Modèle de données

```go
type Paiment struct {
    PaymentId      uuid.UUID
    TenantId       string
    IdempotencyKey string
    OrderRef       string
    Amount         float64
    Currency       string
    Status         Status  // CREATED | SUCCEEDED | FAILED | REQUIRES_ACTION
    CreateAt       time.Time
    NextAction     bool
}
```

## Démarrage

```bash
# Installation des dépendances
go mod tidy

# Lancement du serveur (port 8080)
go run main.go memory.go request.go
```

## Dépendances

- `github.com/gin-gonic/gin` - Framework HTTP
- `github.com/google/uuid` - Génération d'UUID

## Fonctionnalités

- ✅ Création de paiements avec validation
- ✅ Idempotence par `(TenantId, IdempotencyKey)`
- ✅ Gestion des statuts selon le montant
- ✅ Récupération par UUID
- ✅ Stockage en mémoire (non persistant)
