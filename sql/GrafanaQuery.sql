-- For creating new variables
select distinct flat_no from time_series_data order by flat_no asc;

-- Consumption per flat within given time period
select flat_no, total_consumption, timestamp
from time_series_data
where flat_no = 'A001'
order by timestamp DESC, flat_no ASC ;

-- Consumption per month per flat
select time_bucket('1 month', timestamp) as one_month, sum(total_consumption)
from time_series_data
where flat_no = 'A206'
group by one_month
order by one_month desc;

select * from time_series_data;
select metadata from time_series_data where flat_no = 'A005';

select json_object_keys(select metadata from time_series_data where flat_no = 'A005' and timestamp = '2022-08-01 00:00:00.000000 +00:00') as md
select meta.B1 from time_series_data, jsonb_to_record(metadata) as meta(B1 int, B2 int, K1 int) where flat_no = 'A005'
-- truncate table files restart identity;
-- truncate table time_series_data restart identity;

select tsd.flat_no, je.key, je.value, tsd.timestamp
                   from time_series_data tsd,
                        jsonb_each(tsd.metadata) je
                   where tsd.timestamp between '2022-08-01 00:00:00.000000 +00:00' and '2022-08-02 00:00:00.000000 +00:00'
                     and flat_no = 'A005';

select metadata, timestamp from time_series_data
where timestamp between '2022-08-01 00:00:00.000000 +00:00' and '2022-08-02 00:00:00.000000 +00:00'
and flat_no = 'A005';

select jsonb_object_keys(metadata) from time_series_data where flat_no = 'A005' and timestamp = '2022-08-01 00:00:00.000000 +00:00'
create type X as (B1 int, B2 int, K1 int);
select jsonb_populate_record(null::X, metadata) from time_series_data where flat_no = 'A005' and timestamp = '2022-08-01 00:00:00.000000 +00:00'