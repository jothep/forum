使用postgresql数据库

## 初始化数据库
create user -P gwp

createdb -Ogwp -Eutf8 gwp

cd data
psql -f setupforum.sql -U gwp gwp
