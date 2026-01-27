-- 000014_seed_rating_systems.up.sql
-- Seed international rating systems with their ratings

-- ============================================================================
-- MPAA (USA)
-- ============================================================================
INSERT INTO rating_systems (id, code, name, country_codes, sort_order) VALUES
    ('11111111-0001-0001-0001-000000000001', 'mpaa', 'Motion Picture Association', ARRAY['US', 'CA'], 1);

INSERT INTO ratings (id, system_id, code, name, description, min_age, normalized_level, sort_order, is_adult) VALUES
    ('22222222-0001-0001-0001-000000000001', '11111111-0001-0001-0001-000000000001', 'G', 'General Audiences', 'All ages admitted', 0, 0, 1, false),
    ('22222222-0001-0001-0001-000000000002', '11111111-0001-0001-0001-000000000001', 'PG', 'Parental Guidance Suggested', 'Some material may not be suitable for children', 0, 25, 2, false),
    ('22222222-0001-0001-0001-000000000003', '11111111-0001-0001-0001-000000000001', 'PG-13', 'Parents Strongly Cautioned', 'Some material may be inappropriate for children under 13', 13, 50, 3, false),
    ('22222222-0001-0001-0001-000000000004', '11111111-0001-0001-0001-000000000001', 'R', 'Restricted', 'Under 17 requires accompanying parent or adult guardian', 17, 75, 4, false),
    ('22222222-0001-0001-0001-000000000005', '11111111-0001-0001-0001-000000000001', 'NC-17', 'Adults Only', 'No one 17 and under admitted', 18, 90, 5, false),
    ('22222222-0001-0001-0001-000000000006', '11111111-0001-0001-0001-000000000001', 'NR', 'Not Rated', 'Not yet rated or unrated', NULL, 90, 6, false),
    ('22222222-0001-0001-0001-000000000007', '11111111-0001-0001-0001-000000000001', 'XXX', 'Adult', 'Explicit adult content', 18, 100, 7, true);

-- ============================================================================
-- FSK (Germany)
-- ============================================================================
INSERT INTO rating_systems (id, code, name, country_codes, sort_order) VALUES
    ('11111111-0001-0001-0001-000000000002', 'fsk', 'Freiwillige Selbstkontrolle der Filmwirtschaft', ARRAY['DE', 'AT'], 2);

INSERT INTO ratings (id, system_id, code, name, description, min_age, normalized_level, sort_order, is_adult) VALUES
    ('22222222-0002-0001-0001-000000000001', '11111111-0001-0001-0001-000000000002', 'FSK 0', 'Freigegeben ohne Altersbeschränkung', 'Suitable for all ages', 0, 0, 1, false),
    ('22222222-0002-0001-0001-000000000002', '11111111-0001-0001-0001-000000000002', 'FSK 6', 'Freigegeben ab 6 Jahren', 'Suitable for ages 6 and above', 6, 25, 2, false),
    ('22222222-0002-0001-0001-000000000003', '11111111-0001-0001-0001-000000000002', 'FSK 12', 'Freigegeben ab 12 Jahren', 'Suitable for ages 12 and above', 12, 50, 3, false),
    ('22222222-0002-0001-0001-000000000004', '11111111-0001-0001-0001-000000000002', 'FSK 16', 'Freigegeben ab 16 Jahren', 'Suitable for ages 16 and above', 16, 75, 4, false),
    ('22222222-0002-0001-0001-000000000005', '11111111-0001-0001-0001-000000000002', 'FSK 18', 'Keine Jugendfreigabe', 'Not suitable for persons under 18', 18, 90, 5, false),
    ('22222222-0002-0001-0001-000000000006', '11111111-0001-0001-0001-000000000002', 'SPIO/JK', 'Nicht jugendfrei', 'Adult content (SPIO/JK verified)', 18, 100, 6, true);

-- ============================================================================
-- BBFC (UK)
-- ============================================================================
INSERT INTO rating_systems (id, code, name, country_codes, sort_order) VALUES
    ('11111111-0001-0001-0001-000000000003', 'bbfc', 'British Board of Film Classification', ARRAY['GB', 'IE'], 3);

INSERT INTO ratings (id, system_id, code, name, description, min_age, normalized_level, sort_order, is_adult) VALUES
    ('22222222-0003-0001-0001-000000000001', '11111111-0001-0001-0001-000000000003', 'U', 'Universal', 'Suitable for all', 0, 0, 1, false),
    ('22222222-0003-0001-0001-000000000002', '11111111-0001-0001-0001-000000000003', 'PG', 'Parental Guidance', 'General viewing, but some scenes may be unsuitable for young children', 0, 25, 2, false),
    ('22222222-0003-0001-0001-000000000003', '11111111-0001-0001-0001-000000000003', '12A', '12A', 'Suitable for 12 years and over (cinema)', 12, 50, 3, false),
    ('22222222-0003-0001-0001-000000000004', '11111111-0001-0001-0001-000000000003', '12', '12', 'Suitable for 12 years and over (home media)', 12, 50, 4, false),
    ('22222222-0003-0001-0001-000000000005', '11111111-0001-0001-0001-000000000003', '15', '15', 'Suitable only for 15 years and over', 15, 75, 5, false),
    ('22222222-0003-0001-0001-000000000006', '11111111-0001-0001-0001-000000000003', '18', '18', 'Suitable only for adults', 18, 90, 6, false),
    ('22222222-0003-0001-0001-000000000007', '11111111-0001-0001-0001-000000000003', 'R18', 'R18', 'Adult works for licensed premises only', 18, 100, 7, true);

-- ============================================================================
-- PEGI (Europe - Games)
-- ============================================================================
INSERT INTO rating_systems (id, code, name, country_codes, sort_order) VALUES
    ('11111111-0001-0001-0001-000000000004', 'pegi', 'Pan European Game Information', ARRAY['EU'], 4);

INSERT INTO ratings (id, system_id, code, name, description, min_age, normalized_level, sort_order, is_adult) VALUES
    ('22222222-0004-0001-0001-000000000001', '11111111-0001-0001-0001-000000000004', 'PEGI 3', 'PEGI 3', 'Suitable for all ages', 3, 0, 1, false),
    ('22222222-0004-0001-0001-000000000002', '11111111-0001-0001-0001-000000000004', 'PEGI 7', 'PEGI 7', 'Suitable for ages 7 and above', 7, 25, 2, false),
    ('22222222-0004-0001-0001-000000000003', '11111111-0001-0001-0001-000000000004', 'PEGI 12', 'PEGI 12', 'Suitable for ages 12 and above', 12, 50, 3, false),
    ('22222222-0004-0001-0001-000000000004', '11111111-0001-0001-0001-000000000004', 'PEGI 16', 'PEGI 16', 'Suitable for ages 16 and above', 16, 75, 4, false),
    ('22222222-0004-0001-0001-000000000005', '11111111-0001-0001-0001-000000000004', 'PEGI 18', 'PEGI 18', 'Suitable only for adults', 18, 90, 5, false);

-- ============================================================================
-- ACB (Australia)
-- ============================================================================
INSERT INTO rating_systems (id, code, name, country_codes, sort_order) VALUES
    ('11111111-0001-0001-0001-000000000005', 'acb', 'Australian Classification Board', ARRAY['AU', 'NZ'], 5);

INSERT INTO ratings (id, system_id, code, name, description, min_age, normalized_level, sort_order, is_adult) VALUES
    ('22222222-0005-0001-0001-000000000001', '11111111-0001-0001-0001-000000000005', 'G', 'General', 'Suitable for all ages', 0, 0, 1, false),
    ('22222222-0005-0001-0001-000000000002', '11111111-0001-0001-0001-000000000005', 'PG', 'Parental Guidance', 'Parental guidance recommended for children under 15', 0, 25, 2, false),
    ('22222222-0005-0001-0001-000000000003', '11111111-0001-0001-0001-000000000005', 'M', 'Mature', 'Recommended for mature audiences 15 years and over', 15, 50, 3, false),
    ('22222222-0005-0001-0001-000000000004', '11111111-0001-0001-0001-000000000005', 'MA15+', 'Mature Accompanied', 'Restricted to 15 years and over', 15, 75, 4, false),
    ('22222222-0005-0001-0001-000000000005', '11111111-0001-0001-0001-000000000005', 'R18+', 'Restricted', 'Restricted to 18 years and over', 18, 90, 5, false),
    ('22222222-0005-0001-0001-000000000006', '11111111-0001-0001-0001-000000000005', 'X18+', 'Restricted Adult', 'Restricted to adults, explicit content', 18, 100, 6, true);

-- ============================================================================
-- CNC (France)
-- ============================================================================
INSERT INTO rating_systems (id, code, name, country_codes, sort_order) VALUES
    ('11111111-0001-0001-0001-000000000006', 'cnc', 'Centre national du cinéma et de l''image animée', ARRAY['FR', 'BE'], 6);

INSERT INTO ratings (id, system_id, code, name, description, min_age, normalized_level, sort_order, is_adult) VALUES
    ('22222222-0006-0001-0001-000000000001', '11111111-0001-0001-0001-000000000006', 'U', 'Tous publics', 'Suitable for all audiences', 0, 0, 1, false),
    ('22222222-0006-0001-0001-000000000002', '11111111-0001-0001-0001-000000000006', '-10', 'Déconseillé aux moins de 10 ans', 'Not recommended for under 10', 10, 25, 2, false),
    ('22222222-0006-0001-0001-000000000003', '11111111-0001-0001-0001-000000000006', '-12', 'Interdit aux moins de 12 ans', 'Prohibited for under 12', 12, 50, 3, false),
    ('22222222-0006-0001-0001-000000000004', '11111111-0001-0001-0001-000000000006', '-16', 'Interdit aux moins de 16 ans', 'Prohibited for under 16', 16, 75, 4, false),
    ('22222222-0006-0001-0001-000000000005', '11111111-0001-0001-0001-000000000006', '-18', 'Interdit aux moins de 18 ans', 'Prohibited for under 18', 18, 90, 5, false),
    ('22222222-0006-0001-0001-000000000006', '11111111-0001-0001-0001-000000000006', 'X', 'Classé X', 'X-rated pornographic content', 18, 100, 6, true);

-- ============================================================================
-- EIRIN (Japan - Film)
-- ============================================================================
INSERT INTO rating_systems (id, code, name, country_codes, sort_order) VALUES
    ('11111111-0001-0001-0001-000000000007', 'eirin', 'Eiga Rinri Kanri Iinkai', ARRAY['JP'], 7);

INSERT INTO ratings (id, system_id, code, name, description, min_age, normalized_level, sort_order, is_adult) VALUES
    ('22222222-0007-0001-0001-000000000001', '11111111-0001-0001-0001-000000000007', 'G', 'General Audiences', 'Suitable for all ages', 0, 0, 1, false),
    ('22222222-0007-0001-0001-000000000002', '11111111-0001-0001-0001-000000000007', 'PG12', 'Parental Guidance 12', 'Parental guidance required for under 12', 12, 50, 2, false),
    ('22222222-0007-0001-0001-000000000003', '11111111-0001-0001-0001-000000000007', 'R15+', 'Restricted 15', 'Prohibited for under 15', 15, 75, 3, false),
    ('22222222-0007-0001-0001-000000000004', '11111111-0001-0001-0001-000000000007', 'R18+', 'Restricted 18', 'Prohibited for under 18', 18, 90, 4, false);

-- ============================================================================
-- CBFC (India)
-- ============================================================================
INSERT INTO rating_systems (id, code, name, country_codes, sort_order) VALUES
    ('11111111-0001-0001-0001-000000000008', 'cbfc', 'Central Board of Film Certification', ARRAY['IN'], 8);

INSERT INTO ratings (id, system_id, code, name, description, min_age, normalized_level, sort_order, is_adult) VALUES
    ('22222222-0008-0001-0001-000000000001', '11111111-0001-0001-0001-000000000008', 'U', 'Unrestricted', 'Suitable for all ages', 0, 0, 1, false),
    ('22222222-0008-0001-0001-000000000002', '11111111-0001-0001-0001-000000000008', 'UA', 'Unrestricted with Caution', 'Parental guidance for children under 12', 12, 50, 2, false),
    ('22222222-0008-0001-0001-000000000003', '11111111-0001-0001-0001-000000000008', 'A', 'Adults Only', 'Restricted to adults', 18, 90, 3, false),
    ('22222222-0008-0001-0001-000000000004', '11111111-0001-0001-0001-000000000008', 'S', 'Restricted to Specialists', 'Restricted to specialized audiences', 18, 90, 4, false);

-- ============================================================================
-- TV-Parental Guidelines (USA - TV)
-- ============================================================================
INSERT INTO rating_systems (id, code, name, country_codes, sort_order) VALUES
    ('11111111-0001-0001-0001-000000000009', 'tv-pg', 'TV Parental Guidelines', ARRAY['US'], 9);

INSERT INTO ratings (id, system_id, code, name, description, min_age, normalized_level, sort_order, is_adult) VALUES
    ('22222222-0009-0001-0001-000000000001', '11111111-0001-0001-0001-000000000009', 'TV-Y', 'All Children', 'Appropriate for all children', 0, 0, 1, false),
    ('22222222-0009-0001-0001-000000000002', '11111111-0001-0001-0001-000000000009', 'TV-Y7', 'Directed to Older Children', 'Designed for children age 7 and above', 7, 25, 2, false),
    ('22222222-0009-0001-0001-000000000003', '11111111-0001-0001-0001-000000000009', 'TV-G', 'General Audience', 'Most parents would find suitable for all ages', 0, 0, 3, false),
    ('22222222-0009-0001-0001-000000000004', '11111111-0001-0001-0001-000000000009', 'TV-PG', 'Parental Guidance', 'May contain some material unsuitable for children', 0, 25, 4, false),
    ('22222222-0009-0001-0001-000000000005', '11111111-0001-0001-0001-000000000009', 'TV-14', 'Parents Strongly Cautioned', 'May be unsuitable for children under 14', 14, 50, 5, false),
    ('22222222-0009-0001-0001-000000000006', '11111111-0001-0001-0001-000000000009', 'TV-MA', 'Mature Audience Only', 'Specifically designed to be viewed by adults', 17, 75, 6, false);

-- ============================================================================
-- Kijkwijzer (Netherlands)
-- ============================================================================
INSERT INTO rating_systems (id, code, name, country_codes, sort_order) VALUES
    ('11111111-0001-0001-0001-000000000010', 'kijkwijzer', 'Kijkwijzer', ARRAY['NL', 'BE'], 10);

INSERT INTO ratings (id, system_id, code, name, description, min_age, normalized_level, sort_order, is_adult) VALUES
    ('22222222-0010-0001-0001-000000000001', '11111111-0001-0001-0001-000000000010', 'AL', 'Alle leeftijden', 'Suitable for all ages', 0, 0, 1, false),
    ('22222222-0010-0001-0001-000000000002', '11111111-0001-0001-0001-000000000010', '6', '6', 'May be harmful to children under 6', 6, 25, 2, false),
    ('22222222-0010-0001-0001-000000000003', '11111111-0001-0001-0001-000000000010', '9', '9', 'May be harmful to children under 9', 9, 35, 3, false),
    ('22222222-0010-0001-0001-000000000004', '11111111-0001-0001-0001-000000000010', '12', '12', 'May be harmful to children under 12', 12, 50, 4, false),
    ('22222222-0010-0001-0001-000000000005', '11111111-0001-0001-0001-000000000010', '14', '14', 'May be harmful to children under 14', 14, 60, 5, false),
    ('22222222-0010-0001-0001-000000000006', '11111111-0001-0001-0001-000000000010', '16', '16', 'May be harmful to children under 16', 16, 75, 6, false),
    ('22222222-0010-0001-0001-000000000007', '11111111-0001-0001-0001-000000000010', '18', '18', 'Not suitable for persons under 18', 18, 90, 7, false);

-- ============================================================================
-- Rating Equivalents (Cross-references)
-- ============================================================================
-- Map equivalent ratings across systems (bidirectional)
-- Example: MPAA G ≈ FSK 0 ≈ BBFC U ≈ ACB G

-- All Ages equivalents
INSERT INTO rating_equivalents (rating_id, equivalent_rating_id) VALUES
    -- MPAA G ↔ FSK 0
    ('22222222-0001-0001-0001-000000000001', '22222222-0002-0001-0001-000000000001'),
    ('22222222-0002-0001-0001-000000000001', '22222222-0001-0001-0001-000000000001'),
    -- MPAA G ↔ BBFC U
    ('22222222-0001-0001-0001-000000000001', '22222222-0003-0001-0001-000000000001'),
    ('22222222-0003-0001-0001-000000000001', '22222222-0001-0001-0001-000000000001'),
    -- FSK 0 ↔ BBFC U
    ('22222222-0002-0001-0001-000000000001', '22222222-0003-0001-0001-000000000001'),
    ('22222222-0003-0001-0001-000000000001', '22222222-0002-0001-0001-000000000001');

-- PG equivalents
INSERT INTO rating_equivalents (rating_id, equivalent_rating_id) VALUES
    -- MPAA PG ↔ FSK 6
    ('22222222-0001-0001-0001-000000000002', '22222222-0002-0001-0001-000000000002'),
    ('22222222-0002-0001-0001-000000000002', '22222222-0001-0001-0001-000000000002'),
    -- MPAA PG ↔ BBFC PG
    ('22222222-0001-0001-0001-000000000002', '22222222-0003-0001-0001-000000000002'),
    ('22222222-0003-0001-0001-000000000002', '22222222-0001-0001-0001-000000000002');

-- 12/13 equivalents
INSERT INTO rating_equivalents (rating_id, equivalent_rating_id) VALUES
    -- MPAA PG-13 ↔ FSK 12
    ('22222222-0001-0001-0001-000000000003', '22222222-0002-0001-0001-000000000003'),
    ('22222222-0002-0001-0001-000000000003', '22222222-0001-0001-0001-000000000003'),
    -- MPAA PG-13 ↔ BBFC 12A
    ('22222222-0001-0001-0001-000000000003', '22222222-0003-0001-0001-000000000003'),
    ('22222222-0003-0001-0001-000000000003', '22222222-0001-0001-0001-000000000003');

-- 16/17 equivalents
INSERT INTO rating_equivalents (rating_id, equivalent_rating_id) VALUES
    -- MPAA R ↔ FSK 16
    ('22222222-0001-0001-0001-000000000004', '22222222-0002-0001-0001-000000000004'),
    ('22222222-0002-0001-0001-000000000004', '22222222-0001-0001-0001-000000000004'),
    -- MPAA R ↔ BBFC 15
    ('22222222-0001-0001-0001-000000000004', '22222222-0003-0001-0001-000000000005'),
    ('22222222-0003-0001-0001-000000000005', '22222222-0001-0001-0001-000000000004');

-- 18 equivalents
INSERT INTO rating_equivalents (rating_id, equivalent_rating_id) VALUES
    -- MPAA NC-17 ↔ FSK 18
    ('22222222-0001-0001-0001-000000000005', '22222222-0002-0001-0001-000000000005'),
    ('22222222-0002-0001-0001-000000000005', '22222222-0001-0001-0001-000000000005'),
    -- MPAA NC-17 ↔ BBFC 18
    ('22222222-0001-0001-0001-000000000005', '22222222-0003-0001-0001-000000000006'),
    ('22222222-0003-0001-0001-000000000006', '22222222-0001-0001-0001-000000000005');

-- Adult (XXX) equivalents
INSERT INTO rating_equivalents (rating_id, equivalent_rating_id) VALUES
    -- MPAA XXX ↔ FSK SPIO/JK
    ('22222222-0001-0001-0001-000000000007', '22222222-0002-0001-0001-000000000006'),
    ('22222222-0002-0001-0001-000000000006', '22222222-0001-0001-0001-000000000007'),
    -- MPAA XXX ↔ BBFC R18
    ('22222222-0001-0001-0001-000000000007', '22222222-0003-0001-0001-000000000007'),
    ('22222222-0003-0001-0001-000000000007', '22222222-0001-0001-0001-000000000007'),
    -- MPAA XXX ↔ ACB X18+
    ('22222222-0001-0001-0001-000000000007', '22222222-0005-0001-0001-000000000006'),
    ('22222222-0005-0001-0001-000000000006', '22222222-0001-0001-0001-000000000007'),
    -- MPAA XXX ↔ CNC X
    ('22222222-0001-0001-0001-000000000007', '22222222-0006-0001-0001-000000000006'),
    ('22222222-0006-0001-0001-000000000006', '22222222-0001-0001-0001-000000000007');
