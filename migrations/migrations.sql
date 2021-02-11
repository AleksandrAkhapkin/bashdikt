/*##########################################
################ ОСНОВНЫЕ ################
##########################################*/

-- //////////////////////////
-- //Таблица пользователей://
-- //////////////////////////
create table IF NOT EXISTS users
(
	user_id serial not null
		constraint users_pk
			primary key,
	user_role varchar(255) default ''::character varying not null,
	email varchar(255) default ''::character varying not null,
	pass varchar(255) default ''::character varying not null,
	last_name varchar(255) default ''::character varying not null,
	first_name varchar(255) default ''::character varying not null,
	middle_name varchar(255) default ''::character varying not null,
	address varchar(25500) default ''::character varying not null,
	created_at timestamp default now() not null,
	updated_at timestamp default now() not null,
	confirm_email boolean default false not null
);

alter table users owner to bashdi;

create unique index users_email_uindex
	on users (email);

create unique index users_user_id_uindex
	on users (user_id);


-- ////////////////////////////
-- //Таблица доп инфо ученик://
-- ////////////////////////////
create table IF NOT EXISTS student_info
(
	level varchar(255) not null,
	user_id integer not null
);

alter table student_info owner to bashdi;

create unique index student_info_user_id_uindex
	on student_info (user_id);

-- /////////////////////////////////
-- //Таблица доп инфо проверяющий://
-- /////////////////////////////////
create table IF NOT EXISTS teacher_info
(
	user_id integer not null,
	info varchar(25500) not null
);

alter table teacher_info owner to bashdi;

create unique index teacher_info_user_id_uindex
	on teacher_info (user_id);


-- ///////////////////////////
-- //Таблица доп инфо админ://
-- ///////////////////////////
create table IF NOT EXISTS organizer_info
(
	user_id integer not null,
	phone varchar(255) not null,
	soc_url varchar(2550) not null,
	count_student varchar(255) not null,
	format_dictation varchar(255) not null,
	add_email text[] default '{}'::text[],
	add_phone text[] default '{}'::text[]
);

alter table organizer_info owner to bashdi;

create unique index organizer_info_phone_uindex
	on organizer_info (phone);

create unique index organizer_info_user_id_uindex
	on organizer_info (user_id);


-- //////////////////////////////////
-- //Таблица прикрепления учеников://
-- //////////////////////////////////
create table IF NOT EXISTS pin_student
(
	pin_id serial not null
		constraint pin_student_pk
			primary key,
	student_id integer,
	teacher_id integer,
	format_dictation varchar(255) default ''::character varying not null,
	time_pin timestamp default now() not null
);

alter table pin_student owner to bashdi;

create unique index pin_student_pin_id_uindex
	on pin_student (pin_id);


-- //////////////////////////////
-- //Таблица кодов авторизации://
-- //////////////////////////////
create table IF NOT EXISTS auth_code
(
	email varchar(255) not null,
	code varchar(255) not null
);

alter table auth_code owner to bashdi;

create unique index auth_code_code_uindex
	on auth_code (code);

create unique index auth_code_email_uindex
	on auth_code (email);


-- //////////////////////
-- //Таблица диктантов://
-- //////////////////////
create table IF NOT EXISTS dictation
(
	id serial not null
		constraint dictation_pk
			primary key,
	user_id integer not null,
	rating integer default 0 not null,
	text text default ''::text not null,
	status varchar(255) default ''::character varying not null,
	online boolean default false not null,
	send_cert boolean default false not null,
	created_at timestamp default now() not null
);

alter table dictation owner to bashdi;

create unique index dictation_id_uindex
	on dictation (id);

create unique index dictation_user_id_uindex
	on dictation (user_id);


-- ////////////////////////////////////
-- //Таблица комментариев к диктанту://
-- ////////////////////////////////////
create table IF NOT EXISTS markers
(
	text text default ''::character varying not null,
	position integer not null,
	student_id integer not null,
	teacher_id integer not null
);

alter table markers owner to bashdi;



-- ######################################
-- ################ ЛОГИ ################
-- ######################################

-- ////////////////////////////////
-- //Таблица логов пользователей://
-- ////////////////////////////////
create table IF NOT EXISTS request_log
(
	id serial not null
		constraint request_log_pk
			primary key,
	dump_request varchar(100000) default ''::character varying not null,
	ip varchar(255) default ''::character varying not null,
	route varchar(2550) default ''::character varying not null,
	created_at timestamp default now() not null
);

alter table request_log owner to bashdi;

create unique index request_log_id_uindex
	on request_log (id);



