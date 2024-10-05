import {h} from "vue";
import {Icon} from "@iconify/vue";

export const renderIcon = (icon: string) => {
    return () => {
        return h(Icon, {
            icon: icon,
            class: 'text-xl'
        }, {
            default: () => h(icon)
        })
    }
}