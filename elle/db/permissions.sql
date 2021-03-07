create user 'thunder'@'localhost' identified by '';

grant insert, update, select on thunder.articles_content to 'thunder'@'localhost';
grant insert, update, select on thunder.articles_info to 'thunder'@'localhost';

