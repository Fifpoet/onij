import { ref } from 'vue'
import { getNetease } from '@/util'
import { SearchResponse } from '@/api/netease/result.ts'

/** 顶栏全屏搜索与搜索页共用 */
const searchValue = ref('')

export function useSearch() {
    // 多个类型使用默认分页参数, 单个类型可指定
    const handleSearch = async (tps: number[], offset: number, limit: number) => {
        if (!searchValue.value.trim()) return

        const promises = tps.map((typ) => {
            let li = limit, of = offset
            if (tps.length > 1) {
                switch (typ) {
                    case 1:
                        li = 24;
                        break;
                    case 10:
                        li = 6;
                        break;
                    case 100:
                        li = 6;
                        break;
                }
            }
            
            return getNetease<SearchResponse>('/search', {
                "keywords": searchValue.value,
                "limit": li,
                "offset": of,
                "type": typ,
            })
        })
        return Promise.all(promises)
    }

    return { searchValue, handleSearch }
}