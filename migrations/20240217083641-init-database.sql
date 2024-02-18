
-- +migrate Up
CREATE TABLE public.users
(
    id serial NOT NULL,
    active boolean NOT NULL DEFAULT true,
    email text unique NOT NULL,
    password text NOT NULL,
    name text NOT NULL,
    gender varchar(1) NOT NULL,
    picture text,
    benefits jsonb not null default '{"base_swipe": 10}',
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT users_pkey PRIMARY KEY (id)
);
CREATE TABLE public.packages
(
    id serial NOT NULL,
    active boolean NOT NULL DEFAULT true,
    name text NOT NULL,
    active_days int,
    terms text,
    description text,
    benefits jsonb not null,
    price numeric not null,
    type varchar(10) not null,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT packages_pkey PRIMARY KEY (id)
);
CREATE TABLE public.user_packages
(
    id serial NOT NULL,
    user_id int not null,
    package_id int not null,
    expired_at timestamp with time zone,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT user_packages_pkey PRIMARY KEY (id),
    CONSTRAINT users_pack_fkey FOREIGN KEY (user_id)
        REFERENCES public.users (id) MATCH SIMPLE 
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT pack_user_fkey FOREIGN KEY (package_id)
        REFERENCES public.packages (id) MATCH SIMPLE 
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
);
CREATE TABLE public.matches
(
    id serial NOT NULL,
    user_id int not null,
    liked boolean not null,
    user_match_id int not null,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT matches_pkey PRIMARY KEY (id),
    CONSTRAINT user_match_fkey FOREIGN KEY (user_match_id)
        REFERENCES public.users (id) MATCH SIMPLE 
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
);

INSERT INTO public.packages(id,active,name,active_days,terms,description, benefits, type, price) values (1,true,'Premium User',null,'Only can be bought once, if already exist, you not be able to buy again','Yooo! Enjoy your exclusive premium pass, you can swipe all you like. There is no limit to love people right?', '{"is_premium":true}', 'Badge', 120000);

-- +migrate Down
DROP TABLE public.matches;
DROP TABLE public.user_packages;
DROP TABLE public.users;
DROP TABLE public.packages;