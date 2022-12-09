create table if not exists orders(
    place String,
    orderDate DATETIME
) engine = Log();