create table user(
    id bigint(20) not null auto_increment,
    user_id bigint(20) not null,
    username varchar(64) collate utf8mb4_general_ci not null,
    password varchar(64) collate utf8mb4_general_ci not null,
    email  varchar(64)  collate utf8mb4_general_ci,
    gender  tinyint(4)  not null default '0',
    create_time timestamp null default current_timestamp ,
    update_time timestamp  null default current_timestamp on  update current_timestamp ,
    primary key(id),
    unique key idx_username (username) using btree,
    unique key idx_user_id (user_id) using btree
)engine=InnoDB,default charset=utf8mb4 collate=utf8mb4_general_ci;

#存放社区信息表
drop table if exists 'community';
create table community(
    id int(11) not null auto_increment,
    community_id  int(10) unsigned not null,
    community_name varchar(20) collate utf8mb4_general_ci not null,
    introduction varchar(128) collate utf8mb4_general_ci not null,
    create_time timestamp not null default current_timestamp,
    update_time timestamp not null default current_timestamp on update current_timestamp,
    primary key(id),
    unique key idx_community_id (community_id),
    unique key idx_community_name (community_name)
)engine=InnoDB,default charset=utf8mb4 collate=utf8mb4_general_ci;

insert into community values('1','1','Go','Golang','2016-11-01 08:10:10','2016-11-01 08:10:10');
insert into community values('2','2','LeetCode','刷题刷题刷题','2020-01-01 08:00:00','2020-01-01 08:00:00');
insert into community values('3','3','CS:Go','Rush B。。。','2018-08-07 08:30:00','2018-08-07 08:30:00');
insert into community values('4','4','LOL','欢迎来到英雄联盟!','2016-01-01 08:00:00','2016-01-01 08:00:00');
#帖子信息表
#帖子状态：帖子发表要进行审核等
create table post
(
    id           bigint(20)                               not null auto_increment,
    post_id      bigint(20)                               not null comment '帖子id',
    title        varchar(128) collate utf8mb4_general_ci  not null comment '标题',
    content      varchar(8192) collate utf8mb4_general_ci not null comment '内容',
    community_id bigint(20)                               not null comment '社区ID',
    user_id      bigint(20)                               not null comment '作者id',
    status       tinyint(4)                               not null default '1' comment '帖子状态',
    create_time  timestamp                                null     default current_timestamp,
    update_time  timestamp                                null     default current_timestamp on update current_timestamp comment '更新时间',
    primary key (id),
    unique key inx_post_id (post_id),
    key idx_community_id (community_id),
    key idx_user_id(user_id)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 collate=utf8mb4_general_ci;
