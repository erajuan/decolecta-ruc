#!/bin/bash
DB_CONTAINER=setup-postgres-1
unzip /tmp/PadronRUC_202506.zip -d /tmp/
awk -F ',' 'FNR > 1 {print $1 "|"($4=="EMPRESA INDIVIDUAL DE RESPONSABILIDAD LIMITADA" ? "EIRL":($4=="PERSONA NATURAL SIN EMPRESA" ? "0":$4))","($7=="NO DISPONIBLE" ? "0":$7)","($8=="NO DISPONIBLE" ? "0":$8)","($9=="COMPUTARIZADO" ? "C":($9=="MANUAL" ? "M":$9))","($10=="COMPUTARIZADO" ? "C":($10=="MANUAL" ? "M":$10))","($11=="SIN ACTIVIDAD" ? "0":$11)","$12","$16 }' PadronRUC_202506.csv > /tmp/ruc2.csv
iconv -f iso-8859-1 -t utf-8 /tmp/ruc2.csv -o /tmp/ruc2utf8.csv
sudo docker cp /tmp/ruc2utf8.csv "$DB_CONTAINER":/tmp
sudo docker cp ~/cronjobs/sunat_ruc_extenses_sql.sh "$DB_CONTAINER":/tmp/
sudo docker exec -i "$DB_CONTAINER" /bin/bash /tmp/sunat_ruc_extenses_sql.sh
