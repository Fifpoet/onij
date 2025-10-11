import { ref } from 'vue'
import {getNetease} from "@/util";
import {SearchResponse} from "@/api/netease/result.ts";

export function useSearch(of: number, li: number) {
    const searchValue = ref('')
    const handleSearch = async (tps: number[]) => {
        if (!searchValue.value.trim()) return

        const promises = tps.map((typ) => {
            return getNetease<SearchResponse>('/search', {
                "keywords": searchValue.value,
                "limit": li,
                "offset": of,
                "type": typ,
            })
        })
        return Promise.all(promises)
    }

    return {
        searchValue,
        handleSearch
    }
}