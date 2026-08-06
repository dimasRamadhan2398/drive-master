--
-- PostgreSQL database dump
--

\restrict w7YHHP1IG1W7qYSNIIjMJff1psbJdHhS3LJXBUl7fjOUWYzXvfXWlZbHmj2OWSY

-- Dumped from database version 16.14
-- Dumped by pg_dump version 16.14

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

DROP INDEX public.idx_transactions_status;
DROP INDEX public.idx_transactions_payment_method_id;
DROP INDEX public.idx_transactions_payment_id;
DROP INDEX public.idx_refunds_payment_id;
DROP INDEX public.idx_payments_user_id;
DROP INDEX public.idx_payments_status;
DROP INDEX public.idx_payments_payment_method_id;
DROP INDEX public.idx_payments_order_id;
DROP INDEX public.idx_payments_booking_id;
DROP INDEX public.idx_payment_methods_code;
ALTER TABLE ONLY public.transactions DROP CONSTRAINT transactions_pkey;
ALTER TABLE ONLY public.refunds DROP CONSTRAINT refunds_pkey;
ALTER TABLE ONLY public.payments DROP CONSTRAINT payments_pkey;
ALTER TABLE ONLY public.payment_methods DROP CONSTRAINT payment_methods_pkey;
ALTER TABLE public.payment_methods ALTER COLUMN id DROP DEFAULT;
DROP TABLE public.transactions;
DROP TABLE public.refunds;
DROP TABLE public.payments;
DROP SEQUENCE public.payment_methods_id_seq;
DROP TABLE public.payment_methods;
SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: payment_methods; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.payment_methods (
    id bigint NOT NULL,
    code character varying(50) NOT NULL,
    name character varying(100) NOT NULL,
    description text,
    icon character varying(255),
    is_active boolean DEFAULT true,
    sort_order bigint DEFAULT 0,
    gateway character varying(50),
    min_amount numeric(12,2) DEFAULT 0,
    max_amount numeric(12,2) DEFAULT 0,
    fee_type character varying(20),
    fee_amount numeric(12,2) DEFAULT 0,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.payment_methods OWNER TO postgres;

--
-- Name: payment_methods_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.payment_methods_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.payment_methods_id_seq OWNER TO postgres;

--
-- Name: payment_methods_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.payment_methods_id_seq OWNED BY public.payment_methods.id;


--
-- Name: payments; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.payments (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    order_id character varying(50) NOT NULL,
    booking_id uuid,
    user_id uuid NOT NULL,
    amount numeric(12,2) NOT NULL,
    currency character varying(3) DEFAULT 'IDR'::character varying,
    status character varying(20) DEFAULT 'pending'::character varying,
    payment_method_id bigint,
    gateway character varying(50),
    gateway_order_id character varying(100),
    gateway_payment_url character varying(500),
    va_number character varying(100),
    qr_code_url character varying(500),
    metadata jsonb,
    description text,
    expiry_time timestamp with time zone,
    paid_at timestamp with time zone,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.payments OWNER TO postgres;

--
-- Name: refunds; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.refunds (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    payment_id uuid NOT NULL,
    transaction_id uuid NOT NULL,
    amount numeric(12,2) NOT NULL,
    reason text,
    status character varying(20) DEFAULT 'pending'::character varying,
    gateway_refund_id character varying(100),
    processed_at timestamp with time zone,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.refunds OWNER TO postgres;

--
-- Name: transactions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.transactions (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    payment_id uuid NOT NULL,
    type character varying(20) NOT NULL,
    status character varying(20) DEFAULT 'pending'::character varying,
    amount numeric(12,2) NOT NULL,
    currency character varying(3) DEFAULT 'IDR'::character varying,
    gateway character varying(50),
    gateway_txn_id character varying(100),
    gateway_response jsonb,
    error_code character varying(50),
    error_message text,
    processed_at timestamp with time zone,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    payment_method_id bigint
);


ALTER TABLE public.transactions OWNER TO postgres;

--
-- Name: payment_methods id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payment_methods ALTER COLUMN id SET DEFAULT nextval('public.payment_methods_id_seq'::regclass);


--
-- Data for Name: payment_methods; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.payment_methods (id, code, name, description, icon, is_active, sort_order, gateway, min_amount, max_amount, fee_type, fee_amount, created_at, updated_at) FROM stdin;
1	credit_card	Credit Card	Pay with credit card (Visa, Mastercard, JCB)	\N	t	0	\N	0.00	0.00	\N	0.00	2026-06-26 12:47:46.285456+07	2026-06-26 12:47:46.285456+07
2	debit_card	Debit Card	Pay with debit card	\N	t	0	\N	0.00	0.00	\N	0.00	2026-06-26 12:47:46.297662+07	2026-06-26 12:47:46.297662+07
3	bank_transfer	Bank Transfer	Transfer via bank ATM or mobile banking	\N	t	0	\N	0.00	0.00	\N	0.00	2026-06-26 12:47:46.299599+07	2026-06-26 12:47:46.299599+07
4	ewallet	E-Wallet	Pay with e-wallet (GoPay, OVO, Dana, ShopeePay)	\N	t	0	\N	0.00	0.00	\N	0.00	2026-06-26 12:47:46.313205+07	2026-06-26 12:47:46.313205+07
5	qris	QRIS	Scan QR code to pay	\N	t	0	\N	0.00	0.00	\N	0.00	2026-06-26 12:47:46.316378+07	2026-06-26 12:47:46.316378+07
6	cod	Cash on Delivery	Pay when you receive the service	\N	t	0	\N	0.00	0.00	\N	0.00	2026-06-26 12:47:46.319658+07	2026-06-26 12:47:46.319658+07
7	virtual_account	Virtual Account	Pay via virtual account number	\N	t	0	\N	0.00	0.00	\N	0.00	2026-06-26 12:47:46.324232+07	2026-06-26 12:47:46.324232+07
\.


--
-- Data for Name: payments; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.payments (id, order_id, booking_id, user_id, amount, currency, status, payment_method_id, gateway, gateway_order_id, gateway_payment_url, va_number, qr_code_url, metadata, description, expiry_time, paid_at, created_at, updated_at) FROM stdin;
b873f687-463a-48d0-a08f-dad3800eeed8	ORD-20260715075605-499257196	b1c2a1da-b04c-4b1b-bc28-eab853854204	54f2b8ea-d33e-401d-a8e2-ee9705e76d9c	2700000.00	IDR	pending	3	pakasir	ORD-20260715075605-499257196	https://app.pakasir.com/pay/drive-master-indonesia/2700000?order_id=ORD-20260715075605-499257196&redirect=https%3A%2F%2Fdrivemaster.id%2Fauth%2Fpayment-status%3ForderId%3DORD-20260715075605-499257196			{"package_id": "11111111-1111-1111-1111-111111111201", "package_name": "8x", "payment_method": "bank_transfer"}	Pembelian 8x	2026-07-16 07:56:05.732678+07	\N	2026-07-15 07:56:05.732678+07	2026-07-15 07:56:05.732678+07
f2e5ef53-fcb5-4c18-8c0d-fe1f45d44f8e	ORD-20260717102459-2200082089	62f454d9-70db-4fc4-8dae-4e582b24c752	3c09cb50-3d99-4ff8-a57f-44fd4135d070	2150000.00	IDR	success	3	bypass	ORD-20260717102459-2200082089				{"package_id": "11111111-1111-1111-1111-111111111101", "package_name": "6x", "payment_method": "bank_transfer"}	Pembelian 6x (Bypassed Success) - Eki Afifah Rahmawati ()	2026-07-18 10:24:59.274567+07	2026-07-17 10:24:59.274553+07	2026-07-17 10:24:59.274567+07	2026-07-17 10:24:59.274568+07
f4962144-4381-4c27-884d-0451d133dd8b	ORD-20260717102508-2510260168	62f454d9-70db-4fc4-8dae-4e582b24c752	3c09cb50-3d99-4ff8-a57f-44fd4135d070	2150000.00	IDR	success	3	bypass	ORD-20260717102508-2510260168				{"package_id": "11111111-1111-1111-1111-111111111101", "package_name": "6x", "payment_method": "bank_transfer"}	Pembelian 6x (Bypassed Success) - Eki Afifah Rahmawati ()	2026-07-18 10:25:08.313827+07	2026-07-17 10:25:08.313824+07	2026-07-17 10:25:08.313827+07	2026-07-17 10:25:08.313827+07
9a977b9d-85a2-4d6a-bf14-7a52a36f2b9d	ORD-20260717103934-4056619485	5d0a2fba-e002-4397-a062-43fd0644205b	54f2b8ea-d33e-401d-a8e2-ee9705e76d9c	2400000.00	IDR	success	3	bypass	ORD-20260717103934-4056619485				{"package_id": "11111111-1111-1111-1111-111111111101", "package_name": "6x", "payment_method": "bank_transfer"}	Pembelian 6x (Bypassed Success) - Addzkiya Nailah ()	2026-07-18 10:39:34.516228+07	2026-07-17 10:39:34.516221+07	2026-07-17 10:39:34.516228+07	2026-07-17 10:39:34.516228+07
bdf1135e-baee-4d5f-8fd6-79185e4a2605	ORD-20260717125723-532541288	5d0a2fba-e002-4397-a062-43fd0644205b	86f5e950-cedb-495b-9528-48a7dffa6919	2400000.00	IDR	pending	3	pakasir	ORD-20260717125723-532541288	https://app.pakasir.com/pay/drive-master-indonesia/2400000?order_id=ORD-20260717125723-532541288&redirect=https%3A%2F%2Fdrivemaster.id%2Fauth%2Fpayment-status%3ForderId%3DORD-20260717125723-532541288			{"package_id": "11111111-1111-1111-1111-111111111101", "package_name": "6x", "payment_method": "bank_transfer"}	Pembelian 6x	2026-07-18 12:57:23.859437+07	\N	2026-07-17 12:57:23.859438+07	2026-07-17 12:57:23.859438+07
e3ef9065-a83d-459c-b240-5e256311add8	ORD-20260717125754-2316667858	5d0a2fba-e002-4397-a062-43fd0644205b	54f2b8ea-d33e-401d-a8e2-ee9705e76d9c	2400000.00	IDR	pending	3	pakasir	ORD-20260717125754-2316667858	https://app.pakasir.com/pay/drive-master-indonesia/2400000?order_id=ORD-20260717125754-2316667858&redirect=https%3A%2F%2Fdrivemaster.id%2Fauth%2Fpayment-status%3ForderId%3DORD-20260717125754-2316667858			{"package_id": "11111111-1111-1111-1111-111111111101", "package_name": "6x", "payment_method": "bank_transfer"}	Pembelian 6x	2026-07-18 12:57:54.439275+07	\N	2026-07-17 12:57:54.439275+07	2026-07-17 12:57:54.439275+07
b92990f9-5b27-41d4-92a0-66ec7638e775	ORD-20260717130058-2425359922	5d0a2fba-e002-4397-a062-43fd0644205b	54f2b8ea-d33e-401d-a8e2-ee9705e76d9c	2400000.00	IDR	pending	3	pakasir	ORD-20260717130058-2425359922	https://app.pakasir.com/pay/drive-master-indonesia/2400000?order_id=ORD-20260717130058-2425359922&redirect=https%3A%2F%2Fdrivemaster.id%2Fauth%2Fpayment-status%3ForderId%3DORD-20260717130058-2425359922			{"package_id": "11111111-1111-1111-1111-111111111101", "package_name": "6x", "payment_method": "bank_transfer"}	Pembelian 6x	2026-07-18 13:00:58.020612+07	\N	2026-07-17 13:00:58.020612+07	2026-07-17 13:00:58.020612+07
6794ca9c-de5a-420b-bd48-9561792a2fda	ORD-20260720101422-2590824459	f5f4326a-50b1-494d-af01-7b31bfbbde2d	3c09cb50-3d99-4ff8-a57f-44fd4135d070	2150000.00	IDR	pending	3	pakasir	ORD-20260720101422-2590824459	https://app.pakasir.com/pay/drive-master-indonesia/2150000?order_id=ORD-20260720101422-2590824459&redirect=https%3A%2F%2Fdrivemaster.id%2Fauth%2Fpayment-status%3ForderId%3DORD-20260720101422-2590824459			{"package_id": "11111111-1111-1111-1111-111111111101", "package_name": "6x", "payment_method": "bank_transfer"}	Pembelian 6x	2026-07-21 10:14:22.850675+07	\N	2026-07-20 10:14:22.850675+07	2026-07-20 10:14:22.850675+07
0efa01d5-6b6b-4f19-b280-0ef3e4222741	ORD-20260720101519-4251413789	f5f4326a-50b1-494d-af01-7b31bfbbde2d	3c09cb50-3d99-4ff8-a57f-44fd4135d070	2150000.00	IDR	pending	3	pakasir	ORD-20260720101519-4251413789	https://app.pakasir.com/pay/drive-master-indonesia/2150000?order_id=ORD-20260720101519-4251413789&redirect=https%3A%2F%2Fdrivemaster.id%2Fauth%2Fpayment-status%3ForderId%3DORD-20260720101519-4251413789			{"package_id": "11111111-1111-1111-1111-111111111101", "package_name": "6x", "payment_method": "bank_transfer"}	Pembelian 6x	2026-07-21 10:15:19.965612+07	\N	2026-07-20 10:15:19.965612+07	2026-07-20 10:15:19.965612+07
c29e9237-d2ad-49c7-8908-795b809506d0	ORD-20260720101654-605705068	f5f4326a-50b1-494d-af01-7b31bfbbde2d	3c09cb50-3d99-4ff8-a57f-44fd4135d070	2150000.00	IDR	pending	5	pakasir	ORD-20260720101654-605705068	https://app.pakasir.com/pay/drive-master-indonesia/2150000?order_id=ORD-20260720101654-605705068&redirect=https%3A%2F%2Fdrivemaster.id%2Fauth%2Fpayment-status%3ForderId%3DORD-20260720101654-605705068			{"package_id": "11111111-1111-1111-1111-111111111101", "package_name": "6x", "payment_method": "qris"}	Pembelian 6x	2026-07-21 10:16:54.709889+07	\N	2026-07-20 10:16:54.70989+07	2026-07-20 10:16:54.70989+07
1cedcda6-8e44-40ce-a78c-5cfef20c6468	ORD-20260721085334-4263799503	6fe72be4-8563-45e6-9f91-15d026a7a9a7	3c09cb50-3d99-4ff8-a57f-44fd4135d070	2150000.00	IDR	pending	3	pakasir	ORD-20260721085334-4263799503	https://app.pakasir.com/pay/drive-master-indonesia/2150000?order_id=ORD-20260721085334-4263799503&redirect=https%3A%2F%2Fdrivemaster.id%2Fauth%2Fpayment-status%3ForderId%3DORD-20260721085334-4263799503			{"package_id": "11111111-1111-1111-1111-111111111101", "package_name": "6x", "payment_method": "bank_transfer"}	Pembelian 6x	2026-07-22 08:53:34.052861+07	\N	2026-07-21 08:53:34.052862+07	2026-07-21 08:53:34.052862+07
6f4cf2cb-c249-4a20-9667-a4c804d3edb9	ORD-20260721085528-655415931	6fe72be4-8563-45e6-9f91-15d026a7a9a7	3c09cb50-3d99-4ff8-a57f-44fd4135d070	2150000.00	IDR	pending	3	pakasir	ORD-20260721085528-655415931	https://app.pakasir.com/pay/drive-master-indonesia/2150000?order_id=ORD-20260721085528-655415931&redirect=https%3A%2F%2Fdrivemaster.id%2Fauth%2Fpayment-status%3ForderId%3DORD-20260721085528-655415931			{"package_id": "11111111-1111-1111-1111-111111111101", "package_name": "6x", "payment_method": "bank_transfer"}	Pembelian 6x	2026-07-22 08:55:28.053492+07	\N	2026-07-21 08:55:28.053492+07	2026-07-21 08:55:28.053493+07
32ed66c6-117a-48d8-9d6d-eea69895960c	ORD-20260721085548-3686006527	6fe72be4-8563-45e6-9f91-15d026a7a9a7	3c09cb50-3d99-4ff8-a57f-44fd4135d070	2150000.00	IDR	pending	5	pakasir	ORD-20260721085548-3686006527	https://app.pakasir.com/pay/drive-master-indonesia/2150000?order_id=ORD-20260721085548-3686006527&redirect=https%3A%2F%2Fdrivemaster.id%2Fauth%2Fpayment-status%3ForderId%3DORD-20260721085548-3686006527			{"package_id": "11111111-1111-1111-1111-111111111101", "package_name": "6x", "payment_method": "qris"}	Pembelian 6x	2026-07-22 08:55:48.199003+07	\N	2026-07-21 08:55:48.199003+07	2026-07-21 08:55:48.199003+07
e8a7b991-640c-4163-ba35-7c7dc7e2d834	ORD-20260721085614-2140690090	6fe72be4-8563-45e6-9f91-15d026a7a9a7	3c09cb50-3d99-4ff8-a57f-44fd4135d070	2150000.00	IDR	pending	3	pakasir	ORD-20260721085614-2140690090	https://app.pakasir.com/pay/drive-master-indonesia/2150000?order_id=ORD-20260721085614-2140690090&redirect=https%3A%2F%2Fdrivemaster.id%2Fauth%2Fpayment-status%3ForderId%3DORD-20260721085614-2140690090			{"package_id": "11111111-1111-1111-1111-111111111101", "package_name": "6x", "payment_method": "bank_transfer"}	Pembelian 6x	2026-07-22 08:56:14.640571+07	\N	2026-07-21 08:56:14.640572+07	2026-07-21 08:56:14.640572+07
c026aeb8-e20d-4e1b-995e-e5c5242378d8	ORD-20260721085906-1710019476	6fe72be4-8563-45e6-9f91-15d026a7a9a7	3c09cb50-3d99-4ff8-a57f-44fd4135d070	2150000.00	IDR	pending	3	pakasir	ORD-20260721085906-1710019476	https://app.pakasir.com/pay/drive-master-indonesia/2150000?order_id=ORD-20260721085906-1710019476&redirect=https%3A%2F%2Fdrivemaster.id%2Fauth%2Fpayment-status%3ForderId%3DORD-20260721085906-1710019476			{"package_id": "11111111-1111-1111-1111-111111111101", "package_name": "6x", "payment_method": "bank_transfer"}	Pembelian 6x	2026-07-22 08:59:06.55226+07	\N	2026-07-21 08:59:06.552261+07	2026-07-21 08:59:06.552261+07
668591e6-08ed-412e-8586-e78835b82161	ORD-20260721085939-2343323102	0e927e20-ccf8-4144-8fc4-4191ce83f643	54f2b8ea-d33e-401d-a8e2-ee9705e76d9c	2400000.00	IDR	pending	3	pakasir	ORD-20260721085939-2343323102	https://app.pakasir.com/pay/drive-master-indonesia/2400000?order_id=ORD-20260721085939-2343323102&redirect=https%3A%2F%2Fdrivemaster.id%2Fauth%2Fpayment-status%3ForderId%3DORD-20260721085939-2343323102			{"package_id": "11111111-1111-1111-1111-111111111101", "package_name": "6x", "payment_method": "bank_transfer"}	Pembelian 6x	2026-07-22 08:59:39.026516+07	\N	2026-07-21 08:59:39.026516+07	2026-07-21 08:59:39.026516+07
6bce63b2-40b8-47b0-934c-d04f4e29b2ae	ORD-20260721090208-3033406053	6fe72be4-8563-45e6-9f91-15d026a7a9a7	3c09cb50-3d99-4ff8-a57f-44fd4135d070	2150000.00	IDR	pending	3	pakasir	ORD-20260721090208-3033406053	https://app.pakasir.com/pay/drive-master-indonesia/2150000?order_id=ORD-20260721090208-3033406053&redirect=https%3A%2F%2Fdrivemaster.id%2Fauth%2Fpayment-status%3ForderId%3DORD-20260721090208-3033406053			{"package_id": "11111111-1111-1111-1111-111111111101", "package_name": "6x", "payment_method": "bank_transfer"}	Pembelian 6x	2026-07-22 09:02:08.786777+07	\N	2026-07-21 09:02:08.786777+07	2026-07-21 09:02:08.786777+07
24634e24-4399-4413-a490-295766541adb	ORD-20260722125709-3026400426	3c3c0ead-befa-475f-819a-574a45e39543	3c09cb50-3d99-4ff8-a57f-44fd4135d070	2150000.00	IDR	pending	3	pakasir	ORD-20260722125709-3026400426	https://app.pakasir.com/pay/drive-master-indonesia/2150000?order_id=ORD-20260722125709-3026400426&redirect=https%3A%2F%2Fdrivemaster.id%2Fauth%2Fpayment-status%3ForderId%3DORD-20260722125709-3026400426			{"package_id": "11111111-1111-1111-1111-111111111101", "package_name": "6x", "payment_method": "bank_transfer"}	Pembelian 6x	2026-07-23 12:57:09.247003+07	\N	2026-07-22 12:57:09.247003+07	2026-07-22 12:57:09.247004+07
434034c9-440a-4bfd-91d0-8eacde2a6095	ORD-20260722141402-2240993343	3c3c0ead-befa-475f-819a-574a45e39543	3c09cb50-3d99-4ff8-a57f-44fd4135d070	2150000.00	IDR	pending	3	pakasir	ORD-20260722141402-2240993343	https://app.pakasir.com/pay/drive-master-indonesia/2150000?order_id=ORD-20260722141402-2240993343&redirect=https%3A%2F%2Fdrivemaster.id%2Fauth%2Fpayment-status%3ForderId%3DORD-20260722141402-2240993343			{"package_id": "11111111-1111-1111-1111-111111111101", "package_name": "6x", "payment_method": "bank_transfer"}	Pembelian 6x	2026-07-23 14:14:02.46183+07	\N	2026-07-22 14:14:02.46183+07	2026-07-22 14:14:02.46183+07
92f32567-894c-42f1-b1ea-43e5f2e7c9df	ORD-20260722141530-1290318921	3c3c0ead-befa-475f-819a-574a45e39543	3c09cb50-3d99-4ff8-a57f-44fd4135d070	2150000.00	IDR	success	3	pakasir	ORD-20260722141530-1290318921	https://app.pakasir.com/pay/drive-master-indonesia/2150000?order_id=ORD-20260722141530-1290318921&redirect=https%3A%2F%2Fdrivemaster.id%2Fauth%2Fpayment-status%3ForderId%3DORD-20260722141530-1290318921			{"package_id": "11111111-1111-1111-1111-111111111101", "package_name": "6x", "payment_method": "bank_transfer"}	Pembelian 6x	2026-07-23 14:15:30.690232+07	2026-07-22 14:19:36.836918+07	2026-07-22 14:15:30.690232+07	2026-07-22 14:19:36.83966+07
8a189b6f-4603-4f39-b516-e352494d67bf	ORD-20260722142425-71052798	3c3c0ead-befa-475f-819a-574a45e39543	3c09cb50-3d99-4ff8-a57f-44fd4135d070	2150000.00	IDR	pending	3	pakasir	ORD-20260722142425-71052798	https://app.pakasir.com/pay/drive-master-indonesia/2150000?order_id=ORD-20260722142425-71052798&redirect=https%3A%2F%2Fdrivemaster.id%2Fauth%2Fpayment-status%3ForderId%3DORD-20260722142425-71052798			{"package_id": "11111111-1111-1111-1111-111111111101", "package_name": "6x", "payment_method": "bank_transfer"}	Pembelian 6x	2026-07-23 14:24:25.657201+07	\N	2026-07-22 14:24:25.657202+07	2026-07-22 14:24:25.657202+07
014aa784-5a5f-4039-ae57-7587539025ed	ORD-20260729113710-2635739770	dd8a6ee9-e2d5-4e2f-820f-256a843d429d	447c7653-04dc-4550-a703-6c81847b1dd6	2150000.00	IDR	pending	3	pakasir	ORD-20260729113710-2635739770	https://app.pakasir.com/pay/drive-master-indonesia-dev/2150000?order_id=ORD-20260729113710-2635739770&redirect=https%3A%2F%2Fdrivemaster.id%2Fauth%2Fpayment-status%3ForderId%3DORD-20260729113710-2635739770			{"package_id": "11111111-1111-1111-1111-111111111101", "package_name": "6x", "payment_method": "bank_transfer"}	Pembelian 6x	2026-07-30 11:37:10.739856+07	\N	2026-07-29 11:37:10.739856+07	2026-07-29 11:37:10.739856+07
7c42f2e4-21a8-4942-862f-273ef774aa6e	ORD-20260729114104-1095364244	dd8a6ee9-e2d5-4e2f-820f-256a843d429d	447c7653-04dc-4550-a703-6c81847b1dd6	2150000.00	IDR	success	3	pakasir	ORD-20260729114104-1095364244	https://app.pakasir.com/pay/drive-master-indonesia-dev/2150000?order_id=ORD-20260729114104-1095364244&redirect=https%3A%2F%2Fdrivemaster.id%2Fauth%2Fpayment-status%3ForderId%3DORD-20260729114104-1095364244			{"package_id": "11111111-1111-1111-1111-111111111101", "package_name": "6x", "payment_method": "bank_transfer"}	Pembelian 6x	2026-07-30 11:41:04.07785+07	2026-07-29 11:41:11.017915+07	2026-07-29 11:41:04.07785+07	2026-07-29 11:41:11.019912+07
43351225-5132-44c9-96dd-9eba4fd12d1f	ORD-20260729115024-2230753000	b310062a-1787-4481-ada4-875f3e9ed5b0	447c7653-04dc-4550-a703-6c81847b1dd6	2150000.00	IDR	pending	3	pakasir	ORD-20260729115024-2230753000	https://app.pakasir.com/pay/drive-master-indonesia-dev/2150000?order_id=ORD-20260729115024-2230753000&redirect=https%3A%2F%2Fdrivemaster.id%2Fauth%2Fpayment-status%3ForderId%3DORD-20260729115024-2230753000			{"package_id": "11111111-1111-1111-1111-111111111101", "package_name": "6x", "payment_method": "bank_transfer"}	Pembelian 6x	2026-07-30 11:50:24.21213+07	\N	2026-07-29 11:50:24.21213+07	2026-07-29 11:50:24.212131+07
e4f5d21a-5bc7-4b90-8ffc-574561faf60c	ORD-20260729115134-2799962904	63347228-eec3-41e4-b9e0-8ddedcc47b68	447c7653-04dc-4550-a703-6c81847b1dd6	2150000.00	IDR	pending	3	pakasir	ORD-20260729115134-2799962904	https://app.pakasir.com/pay/drive-master-indonesia-dev/2150000?order_id=ORD-20260729115134-2799962904&redirect=https%3A%2F%2Fdrivemaster.id%2Fauth%2Fpayment-status%3ForderId%3DORD-20260729115134-2799962904			{"package_id": "11111111-1111-1111-1111-111111111101", "package_name": "6x", "payment_method": "bank_transfer"}	Pembelian 6x	2026-07-30 11:51:34.695179+07	\N	2026-07-29 11:51:34.695179+07	2026-07-29 11:51:34.695179+07
0631512c-52b5-46cc-98e1-4bef98444063	ORD-20260729115358-3835750604	82df6b26-8a85-452b-a715-835d5c99f24d	447c7653-04dc-4550-a703-6c81847b1dd6	2150000.00	IDR	pending	3	pakasir	ORD-20260729115358-3835750604	https://app.pakasir.com/pay/drive-master-indonesia-dev/2150000?order_id=ORD-20260729115358-3835750604&redirect=https%3A%2F%2Fdrivemaster.id%2Fauth%2Fpayment-status%3ForderId%3DORD-20260729115358-3835750604			{"package_id": "11111111-1111-1111-1111-111111111101", "package_name": "6x", "payment_method": "bank_transfer"}	Pembelian 6x	2026-07-30 11:53:58.664743+07	\N	2026-07-29 11:53:58.664743+07	2026-07-29 11:53:58.664743+07
81ac532b-a7cd-41a4-84a9-bb35176721cd	ORD-20260729133211-1655804971	edec3989-a2d9-4155-bbc5-3fcf11b67123	a428042e-e617-48df-8688-cb9ffa7f8c32	2400000.00	IDR	pending	3	pakasir	ORD-20260729133211-1655804971	https://app.pakasir.com/pay/drive-master-indonesia-dev/2400000?order_id=ORD-20260729133211-1655804971&redirect=https%3A%2F%2Fdrivemaster.id%2Fauth%2Fpayment-status%3ForderId%3DORD-20260729133211-1655804971			{"package_id": "11111111-1111-1111-1111-111111111101", "package_name": "6x", "payment_method": "bank_transfer"}	Pembelian 6x	2026-07-30 13:32:11.289971+07	\N	2026-07-29 13:32:11.289972+07	2026-07-29 13:32:11.289972+07
5ba7aae3-43c6-47a6-85bd-47b7d3c49344	ORD-20260729133437-1223382186	edec3989-a2d9-4155-bbc5-3fcf11b67123	a428042e-e617-48df-8688-cb9ffa7f8c32	2400000.00	IDR	pending	3	pakasir	ORD-20260729133437-1223382186	https://app.pakasir.com/pay/drive-master-indonesia-dev/2400000?order_id=ORD-20260729133437-1223382186&redirect=https%3A%2F%2Fdrivemaster.id%2Fauth%2Fpayment-status%3ForderId%3DORD-20260729133437-1223382186			{"package_id": "11111111-1111-1111-1111-111111111101", "package_name": "6x", "payment_method": "bank_transfer"}	Pembelian 6x	2026-07-30 13:34:37.753364+07	\N	2026-07-29 13:34:37.753364+07	2026-07-29 13:34:37.753364+07
269de843-b1e1-454d-8828-08ba905170f0	ORD-20260730050337-850939243	edec3989-a2d9-4155-bbc5-3fcf11b67123	a428042e-e617-48df-8688-cb9ffa7f8c32	2400000.00	IDR	pending	3	pakasir	ORD-20260730050337-850939243	https://app.pakasir.com/pay/drive-master-indonesia/2400000?order_id=ORD-20260730050337-850939243&redirect=https%3A%2F%2Fdrivemaster.id%2Fauth%2Fpayment-status%3ForderId%3DORD-20260730050337-850939243			{"package_id": "11111111-1111-1111-1111-111111111101", "package_name": "6x", "payment_method": "bank_transfer"}	Pembelian 6x	2026-07-31 05:03:37.250479+07	\N	2026-07-30 05:03:37.25048+07	2026-07-30 05:03:37.25048+07
8a53d30e-8bcc-4be9-9654-a0066f8631a6	ORD-20260730051902-2950266031	7fdb8006-33bd-4116-a104-88076e236d0f	a428042e-e617-48df-8688-cb9ffa7f8c32	10000.00	IDR	success	3	pakasir	ORD-20260730051902-2950266031	https://app.pakasir.com/pay/drive-master-indonesia/10000?order_id=ORD-20260730051902-2950266031&redirect=https%3A%2F%2Fdrivemaster.id%2Fauth%2Fpayment-status%3ForderId%3DORD-20260730051902-2950266031			{"package_id": "11111111-1111-1111-1111-111111111101", "package_name": "6x", "payment_method": "bank_transfer"}	Pembelian 6x	2026-07-31 05:19:02.070056+07	2026-07-30 05:19:38.094401+07	2026-07-30 05:19:02.070056+07	2026-07-30 05:19:38.09795+07
dc77adf6-65ef-4af4-9c0a-f3d17be2a31a	ORD-20260730131533-4015994279	e61625ec-8bbe-4a44-bf7a-116fdfa82e15	447c7653-04dc-4550-a703-6c81847b1dd6	2950000.00	IDR	pending	3	pakasir	ORD-20260730131533-4015994279	https://app.pakasir.com/pay/drive-master-indonesia/2950000?order_id=ORD-20260730131533-4015994279&redirect=https%3A%2F%2Fdrivemaster.id%2Fauth%2Fpayment-status%3ForderId%3DORD-20260730131533-4015994279			{"package_id": "11111111-1111-1111-1111-111111111101", "package_name": "6x", "payment_method": "bank_transfer"}	Pembelian 6x	2026-07-31 13:15:33.711209+07	\N	2026-07-30 13:15:33.71121+07	2026-07-30 13:15:33.71121+07
0e12f5de-747d-4f07-9e6e-0d89e95dd174	ORD-20260730181500-2700287614	32c239f2-8c95-467d-ac6e-ff31e19a895e	61b9e452-7b55-46d4-bc8f-23a954a7ace8	2150000.00	IDR	pending	3	pakasir	ORD-20260730181500-2700287614	https://app.pakasir.com/pay/drive-master-indonesia/2150000?order_id=ORD-20260730181500-2700287614&redirect=https%3A%2F%2Fdrivemaster.id%2Fauth%2Fpayment-status%3ForderId%3DORD-20260730181500-2700287614			{"package_id": "11111111-1111-1111-1111-111111111101", "package_name": "6x", "payment_method": "bank_transfer"}	Pembelian 6x	2026-07-31 18:15:00.196692+07	\N	2026-07-30 18:15:00.196693+07	2026-07-30 18:15:00.196693+07
999ebb6f-cf03-4a7e-818e-ba1b3918111d	ORD-20260730181602-88141378	32c239f2-8c95-467d-ac6e-ff31e19a895e	61b9e452-7b55-46d4-bc8f-23a954a7ace8	2150000.00	IDR	pending	3	pakasir	ORD-20260730181602-88141378	https://app.pakasir.com/pay/drive-master-indonesia/2150000?order_id=ORD-20260730181602-88141378&redirect=https%3A%2F%2Fdrivemaster.id%2Fauth%2Fpayment-status%3ForderId%3DORD-20260730181602-88141378			{"package_id": "11111111-1111-1111-1111-111111111101", "package_name": "6x", "payment_method": "bank_transfer"}	Pembelian 6x	2026-07-31 18:16:02.231689+07	\N	2026-07-30 18:16:02.23169+07	2026-07-30 18:16:02.23169+07
156dcb97-b1fe-4641-9c65-8e368432543b	ORD-20260730185838-3890653926	a0c14f0b-9d06-4a9d-83b2-2ed8dff7ceca	fbb44c95-79aa-421e-ba3c-b4e7b87668d8	10000.00	IDR	success	3	pakasir	ORD-20260730185838-3890653926	https://app.pakasir.com/pay/drive-master-indonesia/10000?order_id=ORD-20260730185838-3890653926&redirect=https%3A%2F%2Fdrivemaster.id%2Fauth%2Fpayment-status%3ForderId%3DORD-20260730185838-3890653926			{"package_id": "4601ad1b-4dc9-4727-b229-43cdf6fddecb", "package_name": "Paket Percobaan", "payment_method": "bank_transfer"}	Pembelian Paket Percobaan	2026-07-31 18:58:38.916051+07	2026-07-30 19:00:13.830967+07	2026-07-30 18:58:38.916052+07	2026-07-30 19:00:13.833349+07
f0417f57-6d40-4aa1-a02e-26e496d2d024	ORD-20260804102153-2348102718	69b0f270-ff15-4352-ba8d-f7f383260fca	3c09cb50-3d99-4ff8-a57f-44fd4135d070	10000.00	IDR	success	3	pakasir	ORD-20260804102153-2348102718	https://app.pakasir.com/pay/drive-master-indonesia-dev/10000?order_id=ORD-20260804102153-2348102718&redirect=https%3A%2F%2Fdrivemaster.id%2Fauth%2Fpayment-status%3ForderId%3DORD-20260804102153-2348102718			{"package_id": "4601ad1b-4dc9-4727-b229-43cdf6fddecb", "package_name": "Paket Percobaan", "payment_method": "bank_transfer"}	Pembelian Paket Percobaan	2026-08-05 10:21:53.744971+07	2026-08-04 10:23:20.715786+07	2026-08-04 10:21:53.744971+07	2026-08-04 10:23:20.718535+07
c6750315-a3fc-481e-ae11-c9397c03e0a2	ORD-20260804102440-503709898	a723de49-a8e2-4d78-884f-203dd4522910	3c09cb50-3d99-4ff8-a57f-44fd4135d070	10000.00	IDR	success	3	pakasir	ORD-20260804102440-503709898	https://app.pakasir.com/pay/drive-master-indonesia-dev/10000?order_id=ORD-20260804102440-503709898&redirect=https%3A%2F%2Fdrivemaster.id%2Fauth%2Fpayment-status%3ForderId%3DORD-20260804102440-503709898			{"package_id": "4601ad1b-4dc9-4727-b229-43cdf6fddecb", "package_name": "Paket Percobaan", "payment_method": "bank_transfer"}	Pembelian Paket Percobaan	2026-08-05 10:24:40.02127+07	2026-08-04 10:33:04.542857+07	2026-08-04 10:24:40.02127+07	2026-08-04 10:33:04.545253+07
77e14104-e45b-4b04-9d24-1809fe9a8365	ORD-20260804103556-1590963151	9737a24d-0d15-48ed-897f-c0df00d57b31	fead48e2-b1b6-4620-9b18-889e64355698	10000.00	IDR	success	3	pakasir	ORD-20260804103556-1590963151	https://app.pakasir.com/pay/drive-master-indonesia-dev/10000?order_id=ORD-20260804103556-1590963151&redirect=https%3A%2F%2Fdrivemaster.id%2Fauth%2Fpayment-status%3ForderId%3DORD-20260804103556-1590963151			{"package_id": "4601ad1b-4dc9-4727-b229-43cdf6fddecb", "package_name": "Paket Percobaan", "payment_method": "bank_transfer"}	Pembelian Paket Percobaan	2026-08-05 10:35:56.172056+07	2026-08-04 10:36:39.445595+07	2026-08-04 10:35:56.172056+07	2026-08-04 10:36:39.447978+07
f4221037-e5b9-4383-9f3f-06053f8795d4	ORD-20260804105141-2660810074	5ed13eae-af47-486d-98a2-0160e99c5290	3c09cb50-3d99-4ff8-a57f-44fd4135d070	10000.00	IDR	success	3	pakasir	ORD-20260804105141-2660810074	https://app.pakasir.com/pay/drive-master-indonesia-dev/10000?order_id=ORD-20260804105141-2660810074&redirect=https%3A%2F%2Fdrivemaster.id%2Fauth%2Fpayment-status%3ForderId%3DORD-20260804105141-2660810074			{"package_id": "4601ad1b-4dc9-4727-b229-43cdf6fddecb", "package_name": "Paket Percobaan", "payment_method": "bank_transfer"}	Pembelian Paket Percobaan	2026-08-05 10:51:41.042744+07	2026-08-04 10:51:50.813763+07	2026-08-04 10:51:41.042744+07	2026-08-04 10:51:50.81517+07
\.


--
-- Data for Name: refunds; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.refunds (id, payment_id, transaction_id, amount, reason, status, gateway_refund_id, processed_at, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: transactions; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.transactions (id, payment_id, type, status, amount, currency, gateway, gateway_txn_id, gateway_response, error_code, error_message, processed_at, created_at, updated_at, payment_method_id) FROM stdin;
c352d34a-c772-45d7-8d7f-d0c2ac9c09c8	b873f687-463a-48d0-a08f-dad3800eeed8	charge	pending	2700000.00	IDR	pakasir	ORD-20260715075605-499257196	{}			\N	2026-07-15 07:56:05.742408+07	2026-07-15 07:56:05.742408+07	3
193e4c79-a1e3-457a-94c8-cd525ee22854	f2e5ef53-fcb5-4c18-8c0d-fe1f45d44f8e	charge	success	2150000.00	IDR	bypass	ORD-20260717102459-2200082089	{}			2026-07-17 10:24:59.274553+07	2026-07-17 10:24:59.288818+07	2026-07-17 10:24:59.288818+07	3
205449d1-4965-4452-95fd-d4a459416ce2	f4962144-4381-4c27-884d-0451d133dd8b	charge	success	2150000.00	IDR	bypass	ORD-20260717102508-2510260168	{}			2026-07-17 10:25:08.313824+07	2026-07-17 10:25:08.316737+07	2026-07-17 10:25:08.316737+07	3
6f508414-f48c-48e6-a209-c263f16fe989	9a977b9d-85a2-4d6a-bf14-7a52a36f2b9d	charge	success	2400000.00	IDR	bypass	ORD-20260717103934-4056619485	{}			2026-07-17 10:39:34.516221+07	2026-07-17 10:39:34.518379+07	2026-07-17 10:39:34.518379+07	3
bd3b5a28-9c28-43ae-af47-d73ec4182b86	bdf1135e-baee-4d5f-8fd6-79185e4a2605	charge	pending	2400000.00	IDR	pakasir	ORD-20260717125723-532541288	{}			\N	2026-07-17 12:57:23.865329+07	2026-07-17 12:57:23.865329+07	3
8219e32a-195d-4c96-a621-7e05f2a9200a	e3ef9065-a83d-459c-b240-5e256311add8	charge	pending	2400000.00	IDR	pakasir	ORD-20260717125754-2316667858	{}			\N	2026-07-17 12:57:54.441574+07	2026-07-17 12:57:54.441574+07	3
228f1677-8563-4af4-974e-1d62acedf44f	b92990f9-5b27-41d4-92a0-66ec7638e775	charge	pending	2400000.00	IDR	pakasir	ORD-20260717130058-2425359922	{}			\N	2026-07-17 13:00:58.022969+07	2026-07-17 13:00:58.022969+07	3
5370199d-032c-4ae2-b448-743f09f2b9fd	6794ca9c-de5a-420b-bd48-9561792a2fda	charge	pending	2150000.00	IDR	pakasir	ORD-20260720101422-2590824459	{}			\N	2026-07-20 10:14:22.86034+07	2026-07-20 10:14:22.86034+07	3
3870f6b8-b1ae-40cf-ae9c-fc0c9fa3258a	0efa01d5-6b6b-4f19-b280-0ef3e4222741	charge	pending	2150000.00	IDR	pakasir	ORD-20260720101519-4251413789	{}			\N	2026-07-20 10:15:19.970426+07	2026-07-20 10:15:19.970426+07	3
133a9323-1e1d-4579-a86d-cdf359df73a8	c29e9237-d2ad-49c7-8908-795b809506d0	charge	pending	2150000.00	IDR	pakasir	ORD-20260720101654-605705068	{}			\N	2026-07-20 10:16:54.713091+07	2026-07-20 10:16:54.713091+07	5
8ad8a33b-ab69-4a32-ba91-29d485c1219c	1cedcda6-8e44-40ce-a78c-5cfef20c6468	charge	pending	2150000.00	IDR	pakasir	ORD-20260721085334-4263799503	{}			\N	2026-07-21 08:53:34.055907+07	2026-07-21 08:53:34.055907+07	3
a03ade1c-8f9b-4521-b480-95856cae1227	6f4cf2cb-c249-4a20-9667-a4c804d3edb9	charge	pending	2150000.00	IDR	pakasir	ORD-20260721085528-655415931	{}			\N	2026-07-21 08:55:28.055874+07	2026-07-21 08:55:28.055874+07	3
cae46a20-2030-42e5-b1e0-5fded1e759a1	32ed66c6-117a-48d8-9d6d-eea69895960c	charge	pending	2150000.00	IDR	pakasir	ORD-20260721085548-3686006527	{}			\N	2026-07-21 08:55:48.20071+07	2026-07-21 08:55:48.20071+07	5
f6d5d10f-0304-4750-acd8-9a28b195cc35	e8a7b991-640c-4163-ba35-7c7dc7e2d834	charge	pending	2150000.00	IDR	pakasir	ORD-20260721085614-2140690090	{}			\N	2026-07-21 08:56:14.642762+07	2026-07-21 08:56:14.642762+07	3
549395e0-3001-477c-b7c2-636890106bbe	c026aeb8-e20d-4e1b-995e-e5c5242378d8	charge	pending	2150000.00	IDR	pakasir	ORD-20260721085906-1710019476	{}			\N	2026-07-21 08:59:06.554916+07	2026-07-21 08:59:06.554916+07	3
d1357c88-66e4-49a1-a454-220a61dd30b8	668591e6-08ed-412e-8586-e78835b82161	charge	pending	2400000.00	IDR	pakasir	ORD-20260721085939-2343323102	{}			\N	2026-07-21 08:59:39.028795+07	2026-07-21 08:59:39.028795+07	3
5e141f57-0e90-4c4a-9129-00c5a566b526	6bce63b2-40b8-47b0-934c-d04f4e29b2ae	charge	pending	2150000.00	IDR	pakasir	ORD-20260721090208-3033406053	{}			\N	2026-07-21 09:02:08.791738+07	2026-07-21 09:02:08.791738+07	3
089596df-9176-477b-9ebd-51e21fd8c0e0	24634e24-4399-4413-a490-295766541adb	charge	pending	2150000.00	IDR	pakasir	ORD-20260722125709-3026400426	{}			\N	2026-07-22 12:57:09.250034+07	2026-07-22 12:57:09.250034+07	3
3dc85ee6-d096-4cfd-9476-d2f8f7b7c456	434034c9-440a-4bfd-91d0-8eacde2a6095	charge	pending	2150000.00	IDR	pakasir	ORD-20260722141402-2240993343	{}			\N	2026-07-22 14:14:02.464373+07	2026-07-22 14:14:02.464373+07	3
62d9a775-b546-463a-add1-e0e0e7cd730d	92f32567-894c-42f1-b1ea-43e5f2e7c9df	charge	success	2150000.00	IDR	pakasir	ORD-20260722141530-1290318921	{}			2026-07-22 14:19:36.846474+07	2026-07-22 14:15:30.693172+07	2026-07-22 14:19:36.847344+07	3
f28e4e6d-1a37-4410-a8c4-9fc750bf2852	8a189b6f-4603-4f39-b516-e352494d67bf	charge	pending	2150000.00	IDR	pakasir	ORD-20260722142425-71052798	{}			\N	2026-07-22 14:24:25.659929+07	2026-07-22 14:24:25.659929+07	3
147a8078-fdc4-429e-b4e3-b4ed40b5c689	014aa784-5a5f-4039-ae57-7587539025ed	charge	pending	2150000.00	IDR	pakasir	ORD-20260729113710-2635739770	{}			\N	2026-07-29 11:37:10.748199+07	2026-07-29 11:37:10.748199+07	3
c6aa9732-1f9f-4af7-abe8-1b57180c86fa	7c42f2e4-21a8-4942-862f-273ef774aa6e	charge	success	2150000.00	IDR	pakasir	ORD-20260729114104-1095364244	{}			2026-07-29 11:41:11.026064+07	2026-07-29 11:41:04.080357+07	2026-07-29 11:41:11.026917+07	3
52d4e545-6f86-4ee6-8b64-2c095f3dcec3	43351225-5132-44c9-96dd-9eba4fd12d1f	charge	pending	2150000.00	IDR	pakasir	ORD-20260729115024-2230753000	{}			\N	2026-07-29 11:50:24.216419+07	2026-07-29 11:50:24.216419+07	3
9da0b4f1-31cf-4e1d-99c6-36cc1bb53bcc	e4f5d21a-5bc7-4b90-8ffc-574561faf60c	charge	pending	2150000.00	IDR	pakasir	ORD-20260729115134-2799962904	{}			\N	2026-07-29 11:51:34.696927+07	2026-07-29 11:51:34.696927+07	3
d887503b-9ae4-41f4-a22c-a1b3048c4252	0631512c-52b5-46cc-98e1-4bef98444063	charge	pending	2150000.00	IDR	pakasir	ORD-20260729115358-3835750604	{}			\N	2026-07-29 11:53:58.666994+07	2026-07-29 11:53:58.666994+07	3
f18ec503-0eab-4c40-a58f-660a25123203	81ac532b-a7cd-41a4-84a9-bb35176721cd	charge	pending	2400000.00	IDR	pakasir	ORD-20260729133211-1655804971	{}			\N	2026-07-29 13:32:11.293228+07	2026-07-29 13:32:11.293228+07	3
0417ed6b-6743-4a52-aef2-dc9df1d3559b	5ba7aae3-43c6-47a6-85bd-47b7d3c49344	charge	pending	2400000.00	IDR	pakasir	ORD-20260729133437-1223382186	{}			\N	2026-07-29 13:34:37.755562+07	2026-07-29 13:34:37.755562+07	3
141e2d3c-a3a5-4cde-8116-4297ea5009f8	269de843-b1e1-454d-8828-08ba905170f0	charge	pending	2400000.00	IDR	pakasir	ORD-20260730050337-850939243	{}			\N	2026-07-30 05:03:37.259709+07	2026-07-30 05:03:37.259709+07	3
5dbc54f5-6029-46cf-8e7f-e18616eea0c4	c6750315-a3fc-481e-ae11-c9397c03e0a2	charge	success	10000.00	IDR	pakasir	ORD-20260804102440-503709898	{}			2026-08-04 10:33:04.549653+07	2026-08-04 10:24:40.024498+07	2026-08-04 10:33:04.550343+07	3
b9f4066d-25bc-4997-a9a2-1e282f6a1401	8a53d30e-8bcc-4be9-9654-a0066f8631a6	charge	success	10000.00	IDR	pakasir	ORD-20260730051902-2950266031	{}			2026-07-30 05:19:38.10444+07	2026-07-30 05:19:02.073043+07	2026-07-30 05:19:38.105119+07	3
21fad622-9c37-4780-965c-01ee64e1122b	dc77adf6-65ef-4af4-9c0a-f3d17be2a31a	charge	pending	2950000.00	IDR	pakasir	ORD-20260730131533-4015994279	{}			\N	2026-07-30 13:15:33.714716+07	2026-07-30 13:15:33.714716+07	3
831d2bb2-f8f5-4036-b1b2-caa46b160105	0e12f5de-747d-4f07-9e6e-0d89e95dd174	charge	pending	2150000.00	IDR	pakasir	ORD-20260730181500-2700287614	{}			\N	2026-07-30 18:15:00.199513+07	2026-07-30 18:15:00.199513+07	3
77c54cb5-90f3-45aa-bd2a-bf7021d60133	999ebb6f-cf03-4a7e-818e-ba1b3918111d	charge	pending	2150000.00	IDR	pakasir	ORD-20260730181602-88141378	{}			\N	2026-07-30 18:16:02.233568+07	2026-07-30 18:16:02.233568+07	3
f993a4d8-fd82-4561-bd5e-45e8f52f82a2	156dcb97-b1fe-4641-9c65-8e368432543b	charge	success	10000.00	IDR	pakasir	ORD-20260730185838-3890653926	{}			2026-07-30 19:00:13.839279+07	2026-07-30 18:58:38.918935+07	2026-07-30 19:00:13.840907+07	3
ab5d7c90-ff77-480f-a6cd-64b00f7390f9	f0417f57-6d40-4aa1-a02e-26e496d2d024	charge	success	10000.00	IDR	pakasir	ORD-20260804102153-2348102718	{}			2026-08-04 10:23:20.725133+07	2026-08-04 10:21:53.756181+07	2026-08-04 10:23:20.725997+07	3
05af4d30-21aa-4c1d-8bc1-66e0377e7912	77e14104-e45b-4b04-9d24-1809fe9a8365	charge	success	10000.00	IDR	pakasir	ORD-20260804103556-1590963151	{}			2026-08-04 10:36:39.457821+07	2026-08-04 10:35:56.174106+07	2026-08-04 10:36:39.459276+07	3
90c84631-1375-4874-a025-7163d5a29d58	f4221037-e5b9-4383-9f3f-06053f8795d4	charge	success	10000.00	IDR	pakasir	ORD-20260804105141-2660810074	{}			2026-08-04 10:51:50.81867+07	2026-08-04 10:51:41.048818+07	2026-08-04 10:51:50.819383+07	3
\.


--
-- Name: payment_methods_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.payment_methods_id_seq', 7, true);


--
-- Name: payment_methods payment_methods_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payment_methods
    ADD CONSTRAINT payment_methods_pkey PRIMARY KEY (id);


--
-- Name: payments payments_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payments
    ADD CONSTRAINT payments_pkey PRIMARY KEY (id);


--
-- Name: refunds refunds_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.refunds
    ADD CONSTRAINT refunds_pkey PRIMARY KEY (id);


--
-- Name: transactions transactions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT transactions_pkey PRIMARY KEY (id);


--
-- Name: idx_payment_methods_code; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_payment_methods_code ON public.payment_methods USING btree (code);


--
-- Name: idx_payments_booking_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_payments_booking_id ON public.payments USING btree (booking_id);


--
-- Name: idx_payments_order_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_payments_order_id ON public.payments USING btree (order_id);


--
-- Name: idx_payments_payment_method_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_payments_payment_method_id ON public.payments USING btree (payment_method_id);


--
-- Name: idx_payments_status; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_payments_status ON public.payments USING btree (status);


--
-- Name: idx_payments_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_payments_user_id ON public.payments USING btree (user_id);


--
-- Name: idx_refunds_payment_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_refunds_payment_id ON public.refunds USING btree (payment_id);


--
-- Name: idx_transactions_payment_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_transactions_payment_id ON public.transactions USING btree (payment_id);


--
-- Name: idx_transactions_payment_method_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_transactions_payment_method_id ON public.transactions USING btree (payment_method_id);


--
-- Name: idx_transactions_status; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_transactions_status ON public.transactions USING btree (status);


--
-- PostgreSQL database dump complete
--

\unrestrict w7YHHP1IG1W7qYSNIIjMJff1psbJdHhS3LJXBUl7fjOUWYzXvfXWlZbHmj2OWSY

