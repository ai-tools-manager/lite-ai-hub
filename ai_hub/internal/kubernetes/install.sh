#!/bin/bash

# Этот скрипт устанавливает KEDA Core, KEDA HTTP, treefik Add-on в ваш кластер.
set -e

# Определяем корневую директорию проекта относительно местоположения скрипта
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
PROJECT_ROOT="$SCRIPT_DIR/../../.."

# Переходим в корневую директорию проекта
echo "Переход в корневую директорию проекта: $PROJECT_ROOT"
cd "$PROJECT_ROOT"

KEDA_NAMESPACE="keda"
KEDA_CORE_VERSION="2.16.1"
KEDA_HTTP_ADDON_VERSION="0.9.0"

echo "1. Добавление Helm-репозитория KEDA..."
helm repo add keda https://kedacore.github.io/charts
helm repo update

echo "2. Установка KEDA Core (версия $KEDA_CORE_VERSION)..."
helm install keda keda/keda \
  --namespace $KEDA_NAMESPACE \
  --create-namespace \
  --version $KEDA_CORE_VERSION \
  --wait

echo "KEDA Core успешно установлена."

echo "3. Установка KEDA HTTP Add-on (версия $KEDA_HTTP_ADDON_VERSION)..."
helm install keda-http-add-on keda/keda-add-ons-http \
  --namespace $KEDA_NAMESPACE \
  --version $KEDA_HTTP_ADDON_VERSION \
  --wait

echo "KEDA HTTP Add-on успешно установлен."

echo "4. Установка Traefik Ingress Controller..."
helm repo add traefik https://helm.traefik.io/traefik
helm repo update
helm install traefik traefik/traefik \
  --namespace traefik \
  --create-namespace \
  --wait

echo "Traefik успешно установлен."

kubectl apply -f ./ai_hub/internal/kubernetes/ingress.yaml
