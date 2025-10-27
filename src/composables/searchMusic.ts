import {ref} from 'vue'
import {getNetease} from "@/util";
import {SearchResponse} from "@/api/netease/result.ts";

export function useSearch() {
    const searchValue = ref('')
    const handleSearch = async (tps: number[]) => {
        if (!searchValue.value.trim()) return

        const promises = tps.map((typ) => {
            let li = 0
            switch (typ) {
                case 1:
                    li = 20;
                    break;
                case 10:
                    li = 8;
                    break;
                case 100:
                    li = 5;
                    break;
            }
            return getNetease<SearchResponse>('/search', {
                "keywords": searchValue.value,
                "limit": li,
                "offset": 0,
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