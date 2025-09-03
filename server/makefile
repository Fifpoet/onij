# proto生成handler接口
hz update  -idl _proto/onij/all.proto -I _proto  -out_dir=. --handler_by_method=true

# 依赖注入更新
cd inject && wire

# 数据库生成dal层model
devtool gendb --model_dir ./model --psm infra.rds.test --db xh_test --table mistake_problem  --dal_dir=./infra/dal