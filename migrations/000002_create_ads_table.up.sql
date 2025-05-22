CREATE TABLE ads
(
    id              bigserial PRIMARY KEY,
    title           VARCHAR(255) NOT NULL,
    description     TEXT,
    price           FLOAT        NOT NULL,
    vin             VARCHAR(255) NOT NULL,
    is_token_minted BOOLEAN               DEFAULT FALSE,
    brand           VARCHAR(255) NOT NULL,
    model           VARCHAR(255) NOT NULL,
    year_of_release SMALLINT     NOT NULL,
    image_url       VARCHAR(511)          DEFAULT 'storage/default.jpg',
    user_id         BIGINT       NOT NULL DEFAULT 1,
    created_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    category        VARCHAR(20),
    reg_number      VARCHAR(20),
    type            VARCHAR(20),
    color           VARCHAR(20),
    hp              VARCHAR(20),
    full_weight     VARCHAR(20),
    solo_weight     VARCHAR(20),

    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE OR REPLACE FUNCTION update_ads_modified_column()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_ads_updated_at
    BEFORE UPDATE
    ON ads
    FOR EACH ROW
EXECUTE FUNCTION update_ads_modified_column();

INSERT INTO ads (
    title,
    description,
    price,
    vin,
    brand,
    model,
    year_of_release,
    category,
    reg_number,
    type,
    color,
    hp,
    full_weight,
    solo_weight,
    user_id,
    image_url
) VALUES
      ('BMW X5 2019', 'Превосходное состояние, полный комплект документов', 4500000, 'WBAKR010XL0A12345', 'BMW', 'X5', 2019, 'Кроссовер', 'А777АА777', 'B', 'Чёрный', '249', '2140', '1870', 1, 'https://upload.wikimedia.org/wikipedia/commons/thumb/f/f1/2019_BMW_X5_M50d_Automatic_3.0.jpg/960px-2019_BMW_X5_M50d_Automatic_3.0.jpg'),
      ('Москвич 412 1985', 'Раритет в отличном состоянии, музейный экземпляр', 350000, 'XTA21234056789012', 'Москвич', '412', 1985, 'Седан', 'МОС1985', 'B', 'Бежевый', '75', '1160', '980', 1, 'https://a.d-cd.net/f2SKve_vx2O5v8aRNL7zz9Jlrvk-1920.jpg'),
      ('ЗИЛ-130 1990', 'Рабочая лошадка, полный привод', 850000, 'ZIL1301990ABCDEFG', 'ЗИЛ', '130', 1990, 'Грузовик', 'ЗИЛ1990', 'C', 'Синий', '150', '4300', '3800', 1, 'https://lh5.googleusercontent.com/proxy/JH5WeJj5MFCAOzuX0tk0SCyTnZppXK4zh8zkTsmJzm_go-EyI47-fSmQBl8wsU2yGoJdLBiOtvUyU2dGcjqWuEojqBbK22WmVm7gppqT9D8'),
      ('Lada Vesta Sport 2022', 'Новый, пробег 5000 км', 1500000, 'XTAVESTA2022SPORT', 'Lada', 'Vesta Sport', 2022, 'Седан', 'С222СС777', 'B', 'Красный', '145', '1230', '1090', 1, 'https://avatars.mds.yandex.net/get-autoru-vos/5985536/0b48ed0fa99ac9713ee94566ced26633/456x342'),
      ('Volkswagen Golf GTI 2015', 'Немецкое качество, тюнингованая версия', 1200000, 'WVWZZZAUZFW123456', 'Volkswagen', 'Golf GTI', 2015, 'Хэтчбек', 'Х123ХХ123', 'B', 'Серый', '220', '1345', '1205', 1, 'https://upload.wikimedia.org/wikipedia/commons/thumb/a/ad/Volkswagen_Golf_GTi_2015_%2816460156619%29.jpg/2560px-Volkswagen_Golf_GTi_2015_%2816460156619%29.jpg'),
      ('ГАЗ-21 Волга 1972', 'Полная реставрация, участвует в ретро-ралли', 900000, 'GAZ21000000000000', 'ГАЗ', '21 Волга', 1972, 'Седан', 'ВОЛ1972', 'B', 'Чёрный', '70', '1460', '1260', 1, 'https://s.auto.drom.ru/img5/catalog/photos/fullsize/gaz/21_volga/gaz_21_volga_210075.jpg'),
      ('УАЗ Патриот 2020', 'Внедорожник для охоты и рыбалки', 950000, 'XTTUAZ2020PATRIOT', 'УАЗ', 'Патриот', 2020, 'Внедорожник', 'У777УУ77', 'B', 'Зелёный', '149', '2070', '1850', 1, 'https://iat.ru/uploads/origin/models/724967/1.webp'),
      ('Hyundai Solaris 2018', 'Экономичный седан для города', 800000, 'Z94CB41BAHR123456', 'Hyundai', 'Solaris', 2018, 'Седан', 'Н555НН777', 'B', 'Белый', '123', '1170', '1040', 1, 'https://topruscar.ru/assets/images/kt/kt2018_hyundai-solaris_001.jpg'),
      ('KAMAZ 5490 2021', 'Дальнобойный тягач с пробегом 200 000 км', 8500000, 'KAM549012345678901', 'KAMAZ', '5490', 2021, 'Тягач', 'КАМ2021', 'C', 'Оранжевый', '400', '7100', '6500', 1, 'https://korib.ru/wp-content/uploads/2021/12/kamaz-5490-033-87s5-tyagach.jpg'),
      ('Tesla Model S 2023', 'Электромобиль с автопилотом', 7000000, '5YJSA1E2XNF123456', 'Tesla', 'Model S', 2023, 'Седан', 'Т999ТТ777', 'B', 'Белый', '1020', '2241', '2100', 1, 'https://d2q97jj8nilsnk.cloudfront.net/images/2023-10-04-09.56.19.jpg'),
      ('ВАЗ-2101 1980', 'Жигули "копейка" в оригинальном состоянии', 250000, 'XTA21001000000000', 'ВАЗ', '2101', 1980, 'Седан', 'ЖИГ1980', 'B', 'Голубой', '64', '955', '840', 1, 'storage/default.jpg'),
      ('Porsche 911 Carrera 2020', 'Спортивный автомобиль мечты', 12500000, 'WP0ZZZ99ZLS123456', 'Porsche', '911 Carrera', 2020, 'Купе', 'П911РР777', 'B', 'Жёлтый', '385', '1505', '1420', 1, 'storage/default.jpg'),
      ('ГАЗель NEXT 2019', 'Грузовик для малого бизнеса', 1200000, 'X7LGAZELNEXT2019', 'ГАЗ', 'ГАЗель NEXT', 2019, 'Грузовик', 'ГАЗ2019', 'C', 'Белый', '149', '3050', '2800', 1, 'storage/default.jpg'),
      ('Chevrolet Niva 2015', 'Внедорожник для бездорожья', 600000, 'X7LNIVA2015ABCDE', 'Chevrolet', 'Niva', 2015, 'Внедорожник', 'Ш666ШШ777', 'B', 'Серебристый', '80', '1410', '1270', 1, 'storage/default.jpg'),
      ('ЯМЗ-650 2010', 'Самосвал для строительных работ', 2500000, 'YAMZ6500000000000', 'ЯМЗ', '650', 2010, 'Самосвал', 'ЯМЗ2010', 'C', 'Жёлтый', '312', '12700', '11500', 1, 'storage/default.jpg');