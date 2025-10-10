import { ref } from 'vue'
import {getNetease} from "@/util";
import {SearchResponse} from "@/api/netease/result.ts";

export function useSearch(of: number, li: number) {
    const searchValue = ref('')
    const handleSearch = async (tps: number[]) => {
        if (!searchValue.value.trim()) return

        const promises = tps.map((typ, index) => {
            console.log(`📝 创建第 ${index + 1} 个 Promise，类型: ${typ}`)
            return getNetease<SearchResponse>('/search', {
                "keywords": searchValue.value,
                "limit": li,
                "offset": of,
                "type": typ,
            })
        })

        console.log('✅ promises 创建完成，开始执行 Promise.all...')
        const resp = await Promise.all(promises)
        console.log('🎉 Promise.all 完成:', resp)

        return resp.flatMap(r => r.result)
    }

    return {
        searchValue,
        handleSearch
    }
}