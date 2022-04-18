CREATE TABLE `mes`.`salsitas_weight_control_template` (
  `id_template` INT NOT NULL AUTO_INCREMENT,
  `format` INT NOT NULL,
  `product` INT NOT NULL,
  `error_limit` FLOAT NOT NULL,
  `down_limit` FLOAT NOT NULL,
  `down_limit_warning` FLOAT NOT NULL,
  `target` FLOAT NOT NULL,
  `up_limit_warning` FLOAT NOT NULL,
  `up_limit` FLOAT NOT NULL,
  `empty_doypack` FLOAT NOT NULL,
  `salsitas_weight_control_templatecol` VARCHAR(45) NOT NULL,
  PRIMARY KEY (`id_template`));