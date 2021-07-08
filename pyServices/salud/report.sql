use accesos;

SELECT
    R.registro_fecha,
    Pr.persona_nombre,
    R.registro_tipo_acceso_eid,
    P.persona_nombre,
    R.registro_temperatura,
    R.registro_alerta
FROM
    registros R
        INNER JOIN
    personas Pr ON R.registro_persona_id = Pr.persona_id
        INNER JOIN
    personas P ON R.registro_persona_id = P.persona_id
WHERE
    DATE(R.registro_fecha) BETWEEN '20200420' AND '20200420'



use accesos;

    SELECT
    registro_id,
    CONCAT('assets/images/registros/P10/',
            registro_id,
            '-T.jpg') AS foto_temperatura,
    CONCAT('assets/images/registros/P10/',
            registro_id,
            '-E.jpg') AS foto_epp,
    CONCAT('assets/images/registros/P10/',
            registro_id,
            '-A1.jpg') AS foto_area,
    registro_fecha,
    Pr.persona_nombre,
    R.registro_tipo_acceso_eid,
    GETELEMENTO(registro_tipo_acceso_eid) AS registro_tipo_acceso,
    P.persona_nombre,
    registro_combo1_id,
    GETLISTA(registro_combo1_id) AS resgitro_combo1,
    registro_combo2_id,
    GETLISTA(registro_combo2_id) AS resgitro_combo2,
    R.registro_temperatura,
    R.registro_alerta
FROM
    registros R
        INNER JOIN
    personas Pr ON R.registro_persona_id = Pr.persona_id
        INNER JOIN
    personas P ON R.registro_persona_id = P.persona_id




    use accesos;

  DELIMITER $$
  CREATE FUNCTION `getLista`(`pLista_Id` INT) RETURNS varchar(100) CHARSET latin1
  reads sql data
  deterministic
  BEGIN
  	Declare vNombre varchar(100);

      SELECT lista_nombre into vNombre
        FROM listas
  	 WHERE lista_id = pLista_Id;

      SET vNombre = ifnull(vNombre,'');

  RETURN vNombre;
  END$$
  DELIMITER ;
