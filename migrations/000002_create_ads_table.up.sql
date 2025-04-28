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

INSERT INTO ads (title, description, price, vin, is_token_minted, brand, model, year_of_release, image_url, user_id)
VALUES ('BMW 3 Series for sale', 'This is a listing for a BMW 3 Series in excellent condition.', 20635.92,
        'VIN00013XYZ123456', false, 'BMW', '3 Series', 2009,
        'https://alpinemss.com/cdn/shop/products/11_f6da57d9-6409-4821-9a09-d6348bc48255.png?v=1644869323', 2);
INSERT INTO ads (title, description, price, vin, is_token_minted, brand, model, year_of_release, image_url)
VALUES ('Kia Rio for sale', 'This is a listing for a Kia Rio in excellent condition.', 13403.18, 'VIN00000XYZ123456',
        false, 'Kia', 'Rio', 2006,
        'https://www.avtogermes.ru/images/marks/kia/rio/iv-restajling/colors/saw/kia-rio-solaris-krs-i-belyj-sploshnoj.ebad82697a13a56bcdd243df527b2a6f.png');
INSERT INTO ads (title, description, price, vin, is_token_minted, brand, model, year_of_release, image_url)
VALUES ('Toyota Corolla for sale', 'This is a listing for a Toyota Corolla in excellent condition.', 12615.63,
        'VIN00001XYZ123456', false, 'Toyota', 'Corolla', 2014,
        'https://img-optimize.toyota-europe.com/resize/ccis/680x680/zip/kz/configurationtype/visual-for-grade-selector/product-token/9690d7b0-6ecb-4086-8806-614af98d384c/grade/44961ff8-07c1-45fb-9e52-b1dae327f313/body/74469257-47ae-46eb-859d-2c693e6726ca/fallback/true/padding/50,50,50,50/image-quality/70/day-exterior-4.png');
INSERT INTO ads (title, description, price, vin, is_token_minted, brand, model, year_of_release, image_url)
VALUES ('Honda Civic for sale', 'This is a listing for a Honda Civic in excellent condition.', 21135.89,
        'VIN00002XYZ123456', false, 'Honda', 'Civic', 2009,
        'https://di-honda-enrollment.s3.amazonaws.com/2021/model-pages/civic_hatchback/trims/Hatchback+EX.jpg');
INSERT INTO ads (title, description, price, vin, is_token_minted, brand, model, year_of_release, image_url)
VALUES ('Mercedes C-Class for sale', 'This is a listing for a Mercedes C-Class in excellent condition.', 24747.46,
        'VIN00003XYZ123456', false, 'Mercedes', 'C-Class', 2018,
        'https://di-honda-enrollment.s3.amazonaws.com/2021/model-pages/civic_hatchback/trims/Hatchback+EX.jpg');
INSERT INTO ads (title, description, price, vin, is_token_minted, brand, model, year_of_release, image_url)
VALUES ('Toyota Corolla for sale', 'This is a listing for a Toyota Corolla in excellent condition.', 27712.9,
        'VIN00004XYZ123456', false, 'Toyota', 'Corolla', 2022,
        'https://img-optimize.toyota-europe.com/resize/ccis/680x680/zip/kz/configurationtype/visual-for-grade-selector/product-token/9690d7b0-6ecb-4086-8806-614af98d384c/grade/44961ff8-07c1-45fb-9e52-b1dae327f313/body/74469257-47ae-46eb-859d-2c693e6726ca/fallback/true/padding/50,50,50,50/image-quality/70/day-exterior-4.png');
INSERT INTO ads (title, description, price, vin, is_token_minted, brand, model, year_of_release)
VALUES ('Mercedes C-Class for sale', 'This is a listing for a Mercedes C-Class in excellent condition.', 18762.39,
        'VIN00005XYZ123456', false, 'Mercedes', 'C-Class', 2021);
INSERT INTO ads (title, description, price, vin, is_token_minted, brand, model, year_of_release)
VALUES ('Volkswagen Golf for sale', 'This is a listing for a Volkswagen Golf in excellent condition.', 13754.13,
        'VIN00006XYZ123456', false, 'Volkswagen', 'Golf', 2005);
INSERT INTO ads (title, description, price, vin, is_token_minted, brand, model, year_of_release)
VALUES ('Chevrolet Cruze for sale', 'This is a listing for a Chevrolet Cruze in excellent condition.', 9286.47,
        'VIN00007XYZ123456', false, 'Chevrolet', 'Cruze', 2014);
INSERT INTO ads (title, description, price, vin, is_token_minted, brand, model, year_of_release)
VALUES ('Ford Focus for sale', 'This is a listing for a Ford Focus in excellent condition.', 25214.84,
        'VIN00008XYZ123456', false, 'Ford', 'Focus', 2016);
INSERT INTO ads (title, description, price, vin, is_token_minted, brand, model, year_of_release)
VALUES ('Audi A4 for sale', 'This is a listing for a Audi A4 in excellent condition.', 29292.83, 'VIN00009XYZ123456',
        false, 'Audi', 'A4', 2015);
INSERT INTO ads (title, description, price, vin, is_token_minted, brand, model, year_of_release)
VALUES ('Kia Rio for sale', 'This is a listing for a Kia Rio in excellent condition.', 17202.19, 'VIN00010XYZ123456',
        false, 'Kia', 'Rio', 2013);
INSERT INTO ads (title, description, price, vin, is_token_minted, brand, model, year_of_release)
VALUES ('Chevrolet Cruze for sale', 'This is a listing for a Chevrolet Cruze in excellent condition.', 16939.94,
        'VIN00011XYZ123456', false, 'Chevrolet', 'Cruze', 2015);
INSERT INTO ads (title, description, price, vin, is_token_minted, brand, model, year_of_release)
VALUES ('Toyota Corolla for sale', 'This is a listing for a Toyota Corolla in excellent condition.', 22942.79,
        'VIN00012XYZ123456', false, 'Toyota', 'Corolla', 2013);
INSERT INTO ads (title, description, price, vin, is_token_minted, brand, model, year_of_release)
VALUES ('Audi A4 for sale', 'This is a listing for a Audi A4 in excellent condition.', 19046.87, 'VIN00014XYZ123456',
        false, 'Audi', 'A4', 2013);
INSERT INTO ads (title, description, price, vin, is_token_minted, brand, model, year_of_release)
VALUES ('Ford Focus for sale', 'This is a listing for a Ford Focus in excellent condition.', 26714.24,
        'VIN00015XYZ123456', false, 'Ford', 'Focus', 2005);
INSERT INTO ads (title, description, price, vin, is_token_minted, brand, model, year_of_release)
VALUES ('Chevrolet Cruze for sale', 'This is a listing for a Chevrolet Cruze in excellent condition.', 5506.49,
        'VIN00016XYZ123456', false, 'Chevrolet', 'Cruze', 2010);
INSERT INTO ads (title, description, price, vin, is_token_minted, brand, model, year_of_release)
VALUES ('Kia Rio for sale', 'This is a listing for a Kia Rio in excellent condition.', 15979.46, 'VIN00017XYZ123456',
        false, 'Kia', 'Rio', 2005);
INSERT INTO ads (title, description, price, vin, is_token_minted, brand, model, year_of_release)
VALUES ('Honda Civic for sale', 'This is a listing for a Honda Civic in excellent condition.', 24290.51,
        'VIN00018XYZ123456', false, 'Honda', 'Civic', 2007);
INSERT INTO ads (title, description, price, vin, is_token_minted, brand, model, year_of_release)
VALUES ('Ford Focus for sale', 'This is a listing for a Ford Focus in excellent condition.', 17267.82,
        'VIN00019XYZ123456', false, 'Ford', 'Focus', 2010);
INSERT INTO ads (title, description, price, vin, is_token_minted, brand, model, year_of_release)
VALUES ('Toyota Corolla for sale', 'This is a listing for a Toyota Corolla in excellent condition.', 3825.07,
        'VIN00020XYZ123456', false, 'Toyota', 'Corolla', 2021);
INSERT INTO ads (title, description, price, vin, is_token_minted, brand, model, year_of_release)
VALUES ('Volkswagen Golf for sale', 'This is a listing for a Volkswagen Golf in excellent condition.', 11184.65,
        'VIN00021XYZ123456', false, 'Volkswagen', 'Golf', 2018);
INSERT INTO ads (title, description, price, vin, is_token_minted, brand, model, year_of_release)
VALUES ('BMW 3 Series for sale', 'This is a listing for a BMW 3 Series in excellent condition.', 23065.7,
        'VIN00022XYZ123456', false, 'BMW', '3 Series', 2019);
INSERT INTO ads (title, description, price, vin, is_token_minted, brand, model, year_of_release)
VALUES ('Honda Civic for sale', 'This is a listing for a Honda Civic in excellent condition.', 29151.11,
        'VIN00023XYZ123456', false, 'Honda', 'Civic', 2021);
INSERT INTO ads (title, description, price, vin, is_token_minted, brand, model, year_of_release)
VALUES ('Toyota Corolla for sale', 'This is a listing for a Toyota Corolla in excellent condition.', 3550.67,
        'VIN00024XYZ123456', false, 'Toyota', 'Corolla', 2020);
INSERT INTO ads (title, description, price, vin, is_token_minted, brand, model, year_of_release)
VALUES ('Hyundai Elantra for sale', 'This is a listing for a Hyundai Elantra in excellent condition.', 25486.6,
        'VIN00025XYZ123456', false, 'Hyundai', 'Elantra', 2006);
INSERT INTO ads (title, description, price, vin, is_token_minted, brand, model, year_of_release)
VALUES ('Audi A4 for sale', 'This is a listing for a Audi A4 in excellent condition.', 22008.78, 'VIN00026XYZ123456',
        false, 'Audi', 'A4', 2022);
INSERT INTO ads (title, description, price, vin, is_token_minted, brand, model, year_of_release)
VALUES ('Audi A4 for sale', 'This is a listing for a Audi A4 in excellent condition.', 14986.0, 'VIN00027XYZ123456',
        false, 'Audi', 'A4', 2012);
INSERT INTO ads (title, description, price, vin, is_token_minted, brand, model, year_of_release)
VALUES ('Ford Focus for sale', 'This is a listing for a Ford Focus in excellent condition.', 21839.04,
        'VIN00028XYZ123456', false, 'Ford', 'Focus', 2011);
INSERT INTO ads (title, description, price, vin, is_token_minted, brand, model, year_of_release)
VALUES ('Hyundai Elantra for sale', 'This is a listing for a Hyundai Elantra in excellent condition.', 29539.79,
        'VIN00029XYZ123456', false, 'Hyundai', 'Elantra', 2009);
INSERT INTO ads (title, description, price, vin, is_token_minted, brand, model, year_of_release)
VALUES ('Volkswagen Golf for sale', 'This is a listing for a Volkswagen Golf in excellent condition.', 11549.59,
        'VIN00030XYZ123456', false, 'Volkswagen', 'Golf', 2018);
INSERT INTO ads (title, description, price, vin, is_token_minted, brand, model, year_of_release)
VALUES ('Audi A4 for sale', 'This is a listing for a Audi A4 in excellent condition.', 25910.95, 'VIN00031XYZ123456',
        false, 'Audi', 'A4', 2020);
INSERT INTO ads (title, description, price, vin, is_token_minted, brand, model, year_of_release)
VALUES ('Mercedes C-Class for sale', 'This is a listing for a Mercedes C-Class in excellent condition.', 6155.32,
        'VIN00032XYZ123456', false, 'Mercedes', 'C-Class', 2022);
INSERT INTO ads (title, description, price, vin, is_token_minted, brand, model, year_of_release)
VALUES ('Chevrolet Cruze for sale', 'This is a listing for a Chevrolet Cruze in excellent condition.', 3394.55,
        'VIN00033XYZ123456', false, 'Chevrolet', 'Cruze', 2019);
INSERT INTO ads (title, description, price, vin, is_token_minted, brand, model, year_of_release)
VALUES ('Honda Civic for sale', 'This is a listing for a Honda Civic in excellent condition.', 15290.77,
        'VIN00034XYZ123456', false, 'Honda', 'Civic', 2012);
INSERT INTO ads (title, description, price, vin, is_token_minted, brand, model, year_of_release)
VALUES ('Honda Civic for sale', 'This is a listing for a Honda Civic in excellent condition.', 22121.53,
        'VIN00035XYZ123456', false, 'Honda', 'Civic', 2011);
INSERT INTO ads (title, description, price, vin, is_token_minted, brand, model, year_of_release)
VALUES ('BMW 3 Series for sale', 'This is a listing for a BMW 3 Series in excellent condition.', 13138.29,
        'VIN00036XYZ123456', false, 'BMW', '3 Series', 2009);
INSERT INTO ads (title, description, price, vin, is_token_minted, brand, model, year_of_release)
VALUES ('Volkswagen Golf for sale', 'This is a listing for a Volkswagen Golf in excellent condition.', 8772.55,
        'VIN00037XYZ123456', false, 'Volkswagen', 'Golf', 2013);
INSERT INTO ads (title, description, price, vin, is_token_minted, brand, model, year_of_release)
VALUES ('Honda Civic for sale', 'This is a listing for a Honda Civic in excellent condition.', 27565.2,
        'VIN00038XYZ123456', false, 'Honda', 'Civic', 2019);
INSERT INTO ads (title, description, price, vin, is_token_minted, brand, model, year_of_release)
VALUES ('Audi A4 for sale', 'This is a listing for a Audi A4 in excellent condition.', 24933.9, 'VIN00039XYZ123456',
        false, 'Audi', 'A4', 2016);
INSERT INTO ads (title, description, price, vin, is_token_minted, brand, model, year_of_release)
VALUES ('Chevrolet Cruze for sale', 'This is a listing for a Chevrolet Cruze in excellent condition.', 19980.21,
        'VIN00040XYZ123456', false, 'Chevrolet', 'Cruze', 2012);
INSERT INTO ads (title, description, price, vin, is_token_minted, brand, model, year_of_release)
VALUES ('Kia Rio for sale', 'This is a listing for a Kia Rio in excellent condition.', 24140.46, 'VIN00041XYZ123456',
        false, 'Kia', 'Rio', 2021);
INSERT INTO ads (title, description, price, vin, is_token_minted, brand, model, year_of_release)
VALUES ('BMW 3 Series for sale', 'This is a listing for a BMW 3 Series in excellent condition.', 12878.53,
        'VIN00042XYZ123456', false, 'BMW', '3 Series', 2008);
INSERT INTO ads (title, description, price, vin, is_token_minted, brand, model, year_of_release)
VALUES ('BMW 3 Series for sale', 'This is a listing for a BMW 3 Series in excellent condition.', 26628.97,
        'VIN00043XYZ123456', false, 'BMW', '3 Series', 2005);
INSERT INTO ads (title, description, price, vin, is_token_minted, brand, model, year_of_release)
VALUES ('Honda Civic for sale', 'This is a listing for a Honda Civic in excellent condition.', 22962.47,
        'VIN00044XYZ123456', false, 'Honda', 'Civic', 2018);
INSERT INTO ads (title, description, price, vin, is_token_minted, brand, model, year_of_release)
VALUES ('Ford Focus for sale', 'This is a listing for a Ford Focus in excellent condition.', 9876.59,
        'VIN00045XYZ123456', false, 'Ford', 'Focus', 2013);
INSERT INTO ads (title, description, price, vin, is_token_minted, brand, model, year_of_release)
VALUES ('Mercedes C-Class for sale', 'This is a listing for a Mercedes C-Class in excellent condition.', 12629.59,
        'VIN00046XYZ123456', false, 'Mercedes', 'C-Class', 2012);
INSERT INTO ads (title, description, price, vin, is_token_minted, brand, model, year_of_release)
VALUES ('Audi A4 for sale', 'This is a listing for a Audi A4 in excellent condition.', 9415.29, 'VIN00047XYZ123456',
        false, 'Audi', 'A4', 2013);
INSERT INTO ads (title, description, price, vin, is_token_minted, brand, model, year_of_release)
VALUES ('Toyota Corolla for sale', 'This is a listing for a Toyota Corolla in excellent condition.', 5592.91,
        'VIN00048XYZ123456', false, 'Toyota', 'Corolla', 2021);
INSERT INTO ads (title, description, price, vin, is_token_minted, brand, model, year_of_release)
VALUES ('Hyundai Elantra for sale', 'This is a listing for a Hyundai Elantra in excellent condition.', 13720.98,
        'VIN00049XYZ123456', false, 'Hyundai', 'Elantra', 2017);