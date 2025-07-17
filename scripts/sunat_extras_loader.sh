#!/bin/bash
DB_CONTAINER=apu-postgres-1
FILES=("sunat_buenos_contribuyentes.zip" "sunat_agentes_retencion.zip" "padron_reducido_local_anexo.zip")
# Descargar
echo -e "\nDescargando archivos\n"
for FILE in "${FILES[@]}"; do
    aws s3 cp "s3://your-bucket-name/sunat/$FILE" "/tmp/$FILE"
done
for FILE in /tmp/*.zip; do
    unzip -o "$FILE" -d /tmp/
done
# copy files to /tmp docker
for FILE in "${FILES[@]}"; do
    CSV_FILE="${FILE/.zip/.csv}"
    sudo docker cp "/tmp/$CSV_FILE" "$DB_CONTAINER":/tmp
done
# upload sql importer
sudo docker cp ~/cronjobs/sunat_extras_sql.sh "$DB_CONTAINER":/tmp/
# exec
echo -e "\nImport data into the table ..."
sudo docker exec -i "$DB_CONTAINER" /bin/bash /tmp/sunat_extras_sql.sh
