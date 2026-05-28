#!/bin/sh

# Dossier de destination
DATA_DIR="./data"
mkdir -p "$DATA_DIR"

# Fichier contenant la liste des URLs (une par ligne)
MAPS_FILE="MAPS.txt"

if [ ! -f "$MAPS_FILE" ]; then
    echo "❌ Fichier $MAPS_FILE non trouvé."
    exit 1
fi

echo "⏳ Vérification des cartes à télécharger..."

while IFS= read -r url || [ -n "$url" ]; do
    # Ignorer les lignes vides et les commentaires
    case "$url" in
        \#*) continue ;;
        "") continue ;;
    esac

    # Extraire le nom du fichier depuis l'URL
    filename=$(basename "$url")
    
    if [ ! -f "$DATA_DIR/$filename" ]; then
        echo "🌐 Téléchargement de $filename depuis $url..."
        wget -q --show-progress -O "$DATA_DIR/$filename" "$url"
        if [ $? -eq 0 ]; then
            echo "✅ $filename téléchargé avec succès."
        else
            echo "❌ Erreur lors du téléchargement de $filename."
        fi
    else
        echo "✅ $filename existe déjà, passage au suivant."
    fi
done < "$MAPS_FILE"

echo "🎯 Toutes les cartes sont prêtes."
