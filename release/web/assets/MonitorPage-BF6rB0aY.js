import{d as ie,aa as v,r as $,bK as Le,bL as Rt,bM as zt,bN as $t,bO as Lt,ab as Bt,av as At,a_ as Wt,a5 as Oe,aN as Mt,b3 as Et,F as be,aF as jt,ar as Gt,c as pe,az as Vt,ad as o,ah as p,ac as z,ai as E,aD as Ot,bl as ge,aj as Be,bP as he,ak as Dt,al as De,bx as Ae,aG as Ht,w as fe,o as He,aC as It,aw as K,bQ as Ft,aH as Nt,an as q,ao as le,ap as Ut,bo as Xt,A as Kt,bc as qt,bR as Yt,ba as Jt,ax as me,as as de,bS as _e,bT as Qt,bU as Zt,bV as ea,bW as ta,bX as aa,a as ra,e as oa,f as ce,i as c,n as Q,t as B,g as I,j as xe,y as W,m as na,l as D,I as We,h as ne,_ as ia}from"./index-7GUcku4b.js";import{P as sa}from"./PageContent-BoIAJ2Wd.js";import{A as la,d as ye,f as da}from"./Select-50Os9jyK.js";const ca=Le(".v-x-scroll",{overflow:"auto",scrollbarWidth:"none"},[Le("&::-webkit-scrollbar",{width:0,height:0})]),ba=ie({name:"XScroll",props:{disabled:Boolean,onScroll:Function},setup(){const e=$(null);function n(d){!(d.currentTarget.offsetWidth<d.currentTarget.scrollWidth)||d.deltaY===0||(d.currentTarget.scrollLeft+=d.deltaY+d.deltaX,d.preventDefault())}const b=Rt();return ca.mount({id:"vueuc/x-scroll",head:!0,anchorMetaName:zt,ssr:b}),Object.assign({selfRef:e,handleWheel:n},{scrollTo(...d){var y;(y=e.value)===null||y===void 0||y.scrollTo(...d)}})},render(){return v("div",{ref:"selfRef",onScroll:this.onScroll,onWheel:this.disabled?void 0:this.handleWheel,class:"v-x-scroll"},this.$slots)}});var pa="Expected a function";function Se(e,n,b){var l=!0,d=!0;if(typeof e!="function")throw new TypeError(pa);return Lt(b)&&(l="leading"in b?!!b.leading:l,d="trailing"in b?!!b.trailing:d),$t(e,n,{leading:l,maxWait:n,trailing:d})}const fa={tabFontSizeSmall:"14px",tabFontSizeMedium:"14px",tabFontSizeLarge:"16px",tabGapSmallLine:"36px",tabGapMediumLine:"36px",tabGapLargeLine:"36px",tabGapSmallLineVertical:"8px",tabGapMediumLineVertical:"8px",tabGapLargeLineVertical:"8px",tabPaddingSmallLine:"6px 0",tabPaddingMediumLine:"10px 0",tabPaddingLargeLine:"14px 0",tabPaddingVerticalSmallLine:"6px 12px",tabPaddingVerticalMediumLine:"8px 16px",tabPaddingVerticalLargeLine:"10px 20px",tabGapSmallBar:"36px",tabGapMediumBar:"36px",tabGapLargeBar:"36px",tabGapSmallBarVertical:"8px",tabGapMediumBarVertical:"8px",tabGapLargeBarVertical:"8px",tabPaddingSmallBar:"4px 0",tabPaddingMediumBar:"6px 0",tabPaddingLargeBar:"10px 0",tabPaddingVerticalSmallBar:"6px 12px",tabPaddingVerticalMediumBar:"8px 16px",tabPaddingVerticalLargeBar:"10px 20px",tabGapSmallCard:"4px",tabGapMediumCard:"4px",tabGapLargeCard:"4px",tabGapSmallCardVertical:"4px",tabGapMediumCardVertical:"4px",tabGapLargeCardVertical:"4px",tabPaddingSmallCard:"8px 16px",tabPaddingMediumCard:"10px 20px",tabPaddingLargeCard:"12px 24px",tabPaddingSmallSegment:"4px 0",tabPaddingMediumSegment:"6px 0",tabPaddingLargeSegment:"8px 0",tabPaddingVerticalLargeSegment:"0 8px",tabPaddingVerticalSmallCard:"8px 12px",tabPaddingVerticalMediumCard:"10px 16px",tabPaddingVerticalLargeCard:"12px 20px",tabPaddingVerticalSmallSegment:"0 4px",tabPaddingVerticalMediumSegment:"0 6px",tabGapSmallSegment:"0",tabGapMediumSegment:"0",tabGapLargeSegment:"0",tabGapSmallSegmentVertical:"0",tabGapMediumSegmentVertical:"0",tabGapLargeSegmentVertical:"0",panePaddingSmall:"8px 0 0 0",panePaddingMedium:"12px 0 0 0",panePaddingLarge:"16px 0 0 0",closeSize:"18px",closeIconSize:"14px"};function ua(e){const{textColor2:n,primaryColor:b,textColorDisabled:l,closeIconColor:d,closeIconColorHover:y,closeIconColorPressed:C,closeColorHover:h,closeColorPressed:k,tabColor:w,baseColor:m,dividerColor:L,fontWeight:A,textColor1:x,borderRadius:_,fontSize:P,fontWeightStrong:M}=e;return Object.assign(Object.assign({},fa),{colorSegment:w,tabFontSizeCard:P,tabTextColorLine:x,tabTextColorActiveLine:b,tabTextColorHoverLine:b,tabTextColorDisabledLine:l,tabTextColorSegment:x,tabTextColorActiveSegment:n,tabTextColorHoverSegment:n,tabTextColorDisabledSegment:l,tabTextColorBar:x,tabTextColorActiveBar:b,tabTextColorHoverBar:b,tabTextColorDisabledBar:l,tabTextColorCard:x,tabTextColorHoverCard:x,tabTextColorActiveCard:b,tabTextColorDisabledCard:l,barColor:b,closeIconColor:d,closeIconColorHover:y,closeIconColorPressed:C,closeColorHover:h,closeColorPressed:k,closeBorderRadius:_,tabColor:w,tabColorSegment:m,tabBorderColor:L,tabFontWeightActive:A,tabFontWeight:A,tabBorderRadius:_,paneTextColor:n,fontWeightStrong:M})}const va={common:Bt,self:ua},Pe=At("n-tabs"),Ie={tab:[String,Number,Object,Function],name:{type:[String,Number],required:!0},disabled:Boolean,displayDirective:{type:String,default:"if"},closable:{type:Boolean,default:void 0},tabProps:Object,label:[String,Number,Object,Function]},Me=ie({__TAB_PANE__:!0,name:"TabPane",alias:["TabPanel"],props:Ie,slots:Object,setup(e){const n=Oe(Pe,null);return n||Wt("tab-pane","`n-tab-pane` must be placed inside `n-tabs`."),{style:n.paneStyleRef,class:n.paneClassRef,mergedClsPrefix:n.mergedClsPrefixRef}},render(){return v("div",{class:[`${this.mergedClsPrefix}-tab-pane`,this.class],style:this.style},this.$slots)}}),ga=Object.assign({internalLeftPadded:Boolean,internalAddable:Boolean,internalCreatedByPane:Boolean},Vt(Ie,["displayDirective"])),we=ie({__TAB__:!0,inheritAttrs:!1,name:"Tab",props:ga,setup(e){const{mergedClsPrefixRef:n,valueRef:b,typeRef:l,closableRef:d,tabStyleRef:y,addTabStyleRef:C,tabClassRef:h,addTabClassRef:k,tabChangeIdRef:w,onBeforeLeaveRef:m,triggerRef:L,handleAdd:A,activateTab:x,handleClose:_}=Oe(Pe);return{trigger:L,mergedClosable:pe(()=>{if(e.internalAddable)return!1;const{closable:P}=e;return P===void 0?d.value:P}),style:y,addStyle:C,tabClass:h,addTabClass:k,clsPrefix:n,value:b,type:l,handleClose(P){P.stopPropagation(),!e.disabled&&_(e.name)},activateTab(){if(e.disabled)return;if(e.internalAddable){A();return}const{name:P}=e,M=++w.id;if(P!==b.value){const{value:O}=m;O?Promise.resolve(O(e.name,b.value)).then(j=>{j&&w.id===M&&x(P)}):x(P)}}}},render(){const{internalAddable:e,clsPrefix:n,name:b,disabled:l,label:d,tab:y,value:C,mergedClosable:h,trigger:k,$slots:{default:w}}=this,m=d??y;return v("div",{class:`${n}-tabs-tab-wrapper`},this.internalLeftPadded?v("div",{class:`${n}-tabs-tab-pad`}):null,v("div",Object.assign({key:b,"data-name":b,"data-disabled":l?!0:void 0},Mt({class:[`${n}-tabs-tab`,C===b&&`${n}-tabs-tab--active`,l&&`${n}-tabs-tab--disabled`,h&&`${n}-tabs-tab--closable`,e&&`${n}-tabs-tab--addable`,e?this.addTabClass:this.tabClass],onClick:k==="click"?this.activateTab:void 0,onMouseenter:k==="hover"?this.activateTab:void 0,style:e?this.addStyle:this.style},this.internalCreatedByPane?this.tabProps||{}:this.$attrs)),v("span",{class:`${n}-tabs-tab__label`},e?v(be,null,v("div",{class:`${n}-tabs-tab__height-placeholder`}," "),v(jt,{clsPrefix:n},{default:()=>v(la,null)})):w?w():typeof m=="object"?m:Et(m??b)),h&&this.type==="card"?v(Gt,{clsPrefix:n,class:`${n}-tabs-tab__close`,onClick:this.handleClose,disabled:l}):null))}}),ha=o("tabs",`
 box-sizing: border-box;
 width: 100%;
 display: flex;
 flex-direction: column;
 transition:
 background-color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
`,[p("segment-type",[o("tabs-rail",[z("&.transition-disabled",[o("tabs-capsule",`
 transition: none;
 `)])])]),p("top",[o("tab-pane",`
 padding: var(--n-pane-padding-top) var(--n-pane-padding-right) var(--n-pane-padding-bottom) var(--n-pane-padding-left);
 `)]),p("left",[o("tab-pane",`
 padding: var(--n-pane-padding-right) var(--n-pane-padding-bottom) var(--n-pane-padding-left) var(--n-pane-padding-top);
 `)]),p("left, right",`
 flex-direction: row;
 `,[o("tabs-bar",`
 width: 2px;
 right: 0;
 transition:
 top .2s var(--n-bezier),
 max-height .2s var(--n-bezier),
 background-color .3s var(--n-bezier);
 `),o("tabs-tab",`
 padding: var(--n-tab-padding-vertical); 
 `)]),p("right",`
 flex-direction: row-reverse;
 `,[o("tab-pane",`
 padding: var(--n-pane-padding-left) var(--n-pane-padding-top) var(--n-pane-padding-right) var(--n-pane-padding-bottom);
 `),o("tabs-bar",`
 left: 0;
 `)]),p("bottom",`
 flex-direction: column-reverse;
 justify-content: flex-end;
 `,[o("tab-pane",`
 padding: var(--n-pane-padding-bottom) var(--n-pane-padding-right) var(--n-pane-padding-top) var(--n-pane-padding-left);
 `),o("tabs-bar",`
 top: 0;
 `)]),o("tabs-rail",`
 position: relative;
 padding: 3px;
 border-radius: var(--n-tab-border-radius);
 width: 100%;
 background-color: var(--n-color-segment);
 transition: background-color .3s var(--n-bezier);
 display: flex;
 align-items: center;
 `,[o("tabs-capsule",`
 border-radius: var(--n-tab-border-radius);
 position: absolute;
 pointer-events: none;
 background-color: var(--n-tab-color-segment);
 box-shadow: 0 1px 3px 0 rgba(0, 0, 0, .08);
 transition: transform 0.3s var(--n-bezier);
 `),o("tabs-tab-wrapper",`
 flex-basis: 0;
 flex-grow: 1;
 display: flex;
 align-items: center;
 justify-content: center;
 `,[o("tabs-tab",`
 overflow: hidden;
 border-radius: var(--n-tab-border-radius);
 width: 100%;
 display: flex;
 align-items: center;
 justify-content: center;
 `,[p("active",`
 font-weight: var(--n-font-weight-strong);
 color: var(--n-tab-text-color-active);
 `),z("&:hover",`
 color: var(--n-tab-text-color-hover);
 `)])])]),p("flex",[o("tabs-nav",`
 width: 100%;
 position: relative;
 `,[o("tabs-wrapper",`
 width: 100%;
 `,[o("tabs-tab",`
 margin-right: 0;
 `)])])]),o("tabs-nav",`
 box-sizing: border-box;
 line-height: 1.5;
 display: flex;
 transition: border-color .3s var(--n-bezier);
 `,[E("prefix, suffix",`
 display: flex;
 align-items: center;
 `),E("prefix","padding-right: 16px;"),E("suffix","padding-left: 16px;")]),p("top, bottom",[o("tabs-nav-scroll-wrapper",[z("&::before",`
 top: 0;
 bottom: 0;
 left: 0;
 width: 20px;
 `),z("&::after",`
 top: 0;
 bottom: 0;
 right: 0;
 width: 20px;
 `),p("shadow-start",[z("&::before",`
 box-shadow: inset 10px 0 8px -8px rgba(0, 0, 0, .12);
 `)]),p("shadow-end",[z("&::after",`
 box-shadow: inset -10px 0 8px -8px rgba(0, 0, 0, .12);
 `)])])]),p("left, right",[o("tabs-nav-scroll-content",`
 flex-direction: column;
 `),o("tabs-nav-scroll-wrapper",[z("&::before",`
 top: 0;
 left: 0;
 right: 0;
 height: 20px;
 `),z("&::after",`
 bottom: 0;
 left: 0;
 right: 0;
 height: 20px;
 `),p("shadow-start",[z("&::before",`
 box-shadow: inset 0 10px 8px -8px rgba(0, 0, 0, .12);
 `)]),p("shadow-end",[z("&::after",`
 box-shadow: inset 0 -10px 8px -8px rgba(0, 0, 0, .12);
 `)])])]),o("tabs-nav-scroll-wrapper",`
 flex: 1;
 position: relative;
 overflow: hidden;
 `,[o("tabs-nav-y-scroll",`
 height: 100%;
 width: 100%;
 overflow-y: auto; 
 scrollbar-width: none;
 `,[z("&::-webkit-scrollbar, &::-webkit-scrollbar-track-piece, &::-webkit-scrollbar-thumb",`
 width: 0;
 height: 0;
 display: none;
 `)]),z("&::before, &::after",`
 transition: box-shadow .3s var(--n-bezier);
 pointer-events: none;
 content: "";
 position: absolute;
 z-index: 1;
 `)]),o("tabs-nav-scroll-content",`
 display: flex;
 position: relative;
 min-width: 100%;
 min-height: 100%;
 width: fit-content;
 box-sizing: border-box;
 `),o("tabs-wrapper",`
 display: inline-flex;
 flex-wrap: nowrap;
 position: relative;
 `),o("tabs-tab-wrapper",`
 display: flex;
 flex-wrap: nowrap;
 flex-shrink: 0;
 flex-grow: 0;
 `),o("tabs-tab",`
 cursor: pointer;
 white-space: nowrap;
 flex-wrap: nowrap;
 display: inline-flex;
 align-items: center;
 color: var(--n-tab-text-color);
 font-size: var(--n-tab-font-size);
 background-clip: padding-box;
 padding: var(--n-tab-padding);
 transition:
 box-shadow .3s var(--n-bezier),
 color .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 `,[p("disabled",{cursor:"not-allowed"}),E("close",`
 margin-left: 6px;
 transition:
 background-color .3s var(--n-bezier),
 color .3s var(--n-bezier);
 `),E("label",`
 display: flex;
 align-items: center;
 z-index: 1;
 `)]),o("tabs-bar",`
 position: absolute;
 bottom: 0;
 height: 2px;
 border-radius: 1px;
 background-color: var(--n-bar-color);
 transition:
 left .2s var(--n-bezier),
 max-width .2s var(--n-bezier),
 opacity .3s var(--n-bezier),
 background-color .3s var(--n-bezier);
 `,[z("&.transition-disabled",`
 transition: none;
 `),p("disabled",`
 background-color: var(--n-tab-text-color-disabled)
 `)]),o("tabs-pane-wrapper",`
 position: relative;
 overflow: hidden;
 transition: max-height .2s var(--n-bezier);
 `),o("tab-pane",`
 color: var(--n-pane-text-color);
 width: 100%;
 transition:
 color .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 opacity .2s var(--n-bezier);
 left: 0;
 right: 0;
 top: 0;
 `,[z("&.next-transition-leave-active, &.prev-transition-leave-active, &.next-transition-enter-active, &.prev-transition-enter-active",`
 transition:
 color .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 transform .2s var(--n-bezier),
 opacity .2s var(--n-bezier);
 `),z("&.next-transition-leave-active, &.prev-transition-leave-active",`
 position: absolute;
 `),z("&.next-transition-enter-from, &.prev-transition-leave-to",`
 transform: translateX(32px);
 opacity: 0;
 `),z("&.next-transition-leave-to, &.prev-transition-enter-from",`
 transform: translateX(-32px);
 opacity: 0;
 `),z("&.next-transition-leave-from, &.next-transition-enter-to, &.prev-transition-leave-from, &.prev-transition-enter-to",`
 transform: translateX(0);
 opacity: 1;
 `)]),o("tabs-tab-pad",`
 box-sizing: border-box;
 width: var(--n-tab-gap);
 flex-grow: 0;
 flex-shrink: 0;
 `),p("line-type, bar-type",[o("tabs-tab",`
 font-weight: var(--n-tab-font-weight);
 box-sizing: border-box;
 vertical-align: bottom;
 `,[z("&:hover",{color:"var(--n-tab-text-color-hover)"}),p("active",`
 color: var(--n-tab-text-color-active);
 font-weight: var(--n-tab-font-weight-active);
 `),p("disabled",{color:"var(--n-tab-text-color-disabled)"})])]),o("tabs-nav",[p("line-type",[p("top",[E("prefix, suffix",`
 border-bottom: 1px solid var(--n-tab-border-color);
 `),o("tabs-nav-scroll-content",`
 border-bottom: 1px solid var(--n-tab-border-color);
 `),o("tabs-bar",`
 bottom: -1px;
 `)]),p("left",[E("prefix, suffix",`
 border-right: 1px solid var(--n-tab-border-color);
 `),o("tabs-nav-scroll-content",`
 border-right: 1px solid var(--n-tab-border-color);
 `),o("tabs-bar",`
 right: -1px;
 `)]),p("right",[E("prefix, suffix",`
 border-left: 1px solid var(--n-tab-border-color);
 `),o("tabs-nav-scroll-content",`
 border-left: 1px solid var(--n-tab-border-color);
 `),o("tabs-bar",`
 left: -1px;
 `)]),p("bottom",[E("prefix, suffix",`
 border-top: 1px solid var(--n-tab-border-color);
 `),o("tabs-nav-scroll-content",`
 border-top: 1px solid var(--n-tab-border-color);
 `),o("tabs-bar",`
 top: -1px;
 `)]),E("prefix, suffix",`
 transition: border-color .3s var(--n-bezier);
 `),o("tabs-nav-scroll-content",`
 transition: border-color .3s var(--n-bezier);
 `),o("tabs-bar",`
 border-radius: 0;
 `)]),p("card-type",[E("prefix, suffix",`
 transition: border-color .3s var(--n-bezier);
 `),o("tabs-pad",`
 flex-grow: 1;
 transition: border-color .3s var(--n-bezier);
 `),o("tabs-tab-pad",`
 transition: border-color .3s var(--n-bezier);
 `),o("tabs-tab",`
 font-weight: var(--n-tab-font-weight);
 border: 1px solid var(--n-tab-border-color);
 background-color: var(--n-tab-color);
 box-sizing: border-box;
 position: relative;
 vertical-align: bottom;
 display: flex;
 justify-content: space-between;
 font-size: var(--n-tab-font-size);
 color: var(--n-tab-text-color);
 `,[p("addable",`
 padding-left: 8px;
 padding-right: 8px;
 font-size: 16px;
 justify-content: center;
 `,[E("height-placeholder",`
 width: 0;
 font-size: var(--n-tab-font-size);
 `),Ot("disabled",[z("&:hover",`
 color: var(--n-tab-text-color-hover);
 `)])]),p("closable","padding-right: 8px;"),p("active",`
 background-color: #0000;
 font-weight: var(--n-tab-font-weight-active);
 color: var(--n-tab-text-color-active);
 `),p("disabled","color: var(--n-tab-text-color-disabled);")])]),p("left, right",`
 flex-direction: column; 
 `,[E("prefix, suffix",`
 padding: var(--n-tab-padding-vertical);
 `),o("tabs-wrapper",`
 flex-direction: column;
 `),o("tabs-tab-wrapper",`
 flex-direction: column;
 `,[o("tabs-tab-pad",`
 height: var(--n-tab-gap-vertical);
 width: 100%;
 `)])]),p("top",[p("card-type",[o("tabs-scroll-padding","border-bottom: 1px solid var(--n-tab-border-color);"),E("prefix, suffix",`
 border-bottom: 1px solid var(--n-tab-border-color);
 `),o("tabs-tab",`
 border-top-left-radius: var(--n-tab-border-radius);
 border-top-right-radius: var(--n-tab-border-radius);
 `,[p("active",`
 border-bottom: 1px solid #0000;
 `)]),o("tabs-tab-pad",`
 border-bottom: 1px solid var(--n-tab-border-color);
 `),o("tabs-pad",`
 border-bottom: 1px solid var(--n-tab-border-color);
 `)])]),p("left",[p("card-type",[o("tabs-scroll-padding","border-right: 1px solid var(--n-tab-border-color);"),E("prefix, suffix",`
 border-right: 1px solid var(--n-tab-border-color);
 `),o("tabs-tab",`
 border-top-left-radius: var(--n-tab-border-radius);
 border-bottom-left-radius: var(--n-tab-border-radius);
 `,[p("active",`
 border-right: 1px solid #0000;
 `)]),o("tabs-tab-pad",`
 border-right: 1px solid var(--n-tab-border-color);
 `),o("tabs-pad",`
 border-right: 1px solid var(--n-tab-border-color);
 `)])]),p("right",[p("card-type",[o("tabs-scroll-padding","border-left: 1px solid var(--n-tab-border-color);"),E("prefix, suffix",`
 border-left: 1px solid var(--n-tab-border-color);
 `),o("tabs-tab",`
 border-top-right-radius: var(--n-tab-border-radius);
 border-bottom-right-radius: var(--n-tab-border-radius);
 `,[p("active",`
 border-left: 1px solid #0000;
 `)]),o("tabs-tab-pad",`
 border-left: 1px solid var(--n-tab-border-color);
 `),o("tabs-pad",`
 border-left: 1px solid var(--n-tab-border-color);
 `)])]),p("bottom",[p("card-type",[o("tabs-scroll-padding","border-top: 1px solid var(--n-tab-border-color);"),E("prefix, suffix",`
 border-top: 1px solid var(--n-tab-border-color);
 `),o("tabs-tab",`
 border-bottom-left-radius: var(--n-tab-border-radius);
 border-bottom-right-radius: var(--n-tab-border-radius);
 `,[p("active",`
 border-top: 1px solid #0000;
 `)]),o("tabs-tab-pad",`
 border-top: 1px solid var(--n-tab-border-color);
 `),o("tabs-pad",`
 border-top: 1px solid var(--n-tab-border-color);
 `)])])])]),ma=Object.assign(Object.assign({},De.props),{value:[String,Number],defaultValue:[String,Number],trigger:{type:String,default:"click"},type:{type:String,default:"bar"},closable:Boolean,justifyContent:String,size:{type:String,default:"medium"},placement:{type:String,default:"top"},tabStyle:[String,Object],tabClass:String,addTabStyle:[String,Object],addTabClass:String,barWidth:Number,paneClass:String,paneStyle:[String,Object],paneWrapperClass:String,paneWrapperStyle:[String,Object],addable:[Boolean,Object],tabsPadding:{type:Number,default:0},animated:Boolean,onBeforeLeave:Function,onAdd:Function,"onUpdate:value":[Function,Array],onUpdateValue:[Function,Array],onClose:[Function,Array],labelSize:String,activeName:[String,Number],onActiveNameChange:[Function,Array]}),xa=ie({name:"Tabs",props:ma,slots:Object,setup(e,{slots:n}){var b,l,d,y;const{mergedClsPrefixRef:C,inlineThemeDisabled:h}=Dt(e),k=De("Tabs","-tabs",ha,va,e,C),w=$(null),m=$(null),L=$(null),A=$(null),x=$(null),_=$(null),P=$(!0),M=$(!0),O=Ae(e,["labelSize","size"]),j=Ae(e,["activeName","value"]),Y=$((l=(b=j.value)!==null&&b!==void 0?b:e.defaultValue)!==null&&l!==void 0?l:n.default?(y=(d=ge(n.default())[0])===null||d===void 0?void 0:d.props)===null||y===void 0?void 0:y.name:null),V=Ht(j,Y),S={id:0},G=pe(()=>{if(!(!e.justifyContent||e.type==="card"))return{display:"flex",justifyContent:e.justifyContent}});fe(V,()=>{S.id=0,g(),T()});function H(){var t;const{value:a}=V;return a===null?null:(t=w.value)===null||t===void 0?void 0:t.querySelector(`[data-name="${a}"]`)}function u(t){if(e.type==="card")return;const{value:a}=m;if(!a)return;const s=a.style.opacity==="0";if(t){const f=`${C.value}-tabs-bar--disabled`,{barWidth:R,placement:U}=e;if(t.dataset.disabled==="true"?a.classList.add(f):a.classList.remove(f),["top","bottom"].includes(U)){if(r(["top","maxHeight","height"]),typeof R=="number"&&t.offsetWidth>=R){const X=Math.floor((t.offsetWidth-R)/2)+t.offsetLeft;a.style.left=`${X}px`,a.style.maxWidth=`${R}px`}else a.style.left=`${t.offsetLeft}px`,a.style.maxWidth=`${t.offsetWidth}px`;a.style.width="8192px",s&&(a.style.transition="none"),a.offsetWidth,s&&(a.style.transition="",a.style.opacity="1")}else{if(r(["left","maxWidth","width"]),typeof R=="number"&&t.offsetHeight>=R){const X=Math.floor((t.offsetHeight-R)/2)+t.offsetTop;a.style.top=`${X}px`,a.style.maxHeight=`${R}px`}else a.style.top=`${t.offsetTop}px`,a.style.maxHeight=`${t.offsetHeight}px`;a.style.height="8192px",s&&(a.style.transition="none"),a.offsetHeight,s&&(a.style.transition="",a.style.opacity="1")}}}function i(){if(e.type==="card")return;const{value:t}=m;t&&(t.style.opacity="0")}function r(t){const{value:a}=m;if(a)for(const s of t)a.style[s]=""}function g(){if(e.type==="card")return;const t=H();t?u(t):i()}function T(){var t;const a=(t=x.value)===null||t===void 0?void 0:t.$el;if(!a)return;const s=H();if(!s)return;const{scrollLeft:f,offsetWidth:R}=a,{offsetLeft:U,offsetWidth:X}=s;f>U?a.scrollTo({top:0,left:U,behavior:"smooth"}):U+X>f+R&&a.scrollTo({top:0,left:U+X-R,behavior:"smooth"})}const N=$(null);let Z=0,J=null;function Ne(t){const a=N.value;if(a){Z=t.getBoundingClientRect().height;const s=`${Z}px`,f=()=>{a.style.height=s,a.style.maxHeight=s};J?(f(),J(),J=null):J=f}}function Ue(t){const a=N.value;if(a){const s=t.getBoundingClientRect().height,f=()=>{document.body.offsetHeight,a.style.maxHeight=`${s}px`,a.style.height=`${Math.max(Z,s)}px`};J?(J(),J=null,f()):J=f}}function Xe(){const t=N.value;if(t){t.style.maxHeight="",t.style.height="";const{paneWrapperStyle:a}=e;if(typeof a=="string")t.style.cssText=a;else if(a){const{maxHeight:s,height:f}=a;s!==void 0&&(t.style.maxHeight=s),f!==void 0&&(t.style.height=f)}}}const Te={value:[]},ke=$("next");function Ke(t){const a=V.value;let s="next";for(const f of Te.value){if(f===a)break;if(f===t){s="prev";break}}ke.value=s,qe(t)}function qe(t){const{onActiveNameChange:a,onUpdateValue:s,"onUpdate:value":f}=e;a&&de(a,t),s&&de(s,t),f&&de(f,t),Y.value=t}function Ye(t){const{onClose:a}=e;a&&de(a,t)}function Re(){const{value:t}=m;if(!t)return;const a="transition-disabled";t.classList.add(a),g(),t.classList.remove(a)}const ee=$(null);function ue({transitionDisabled:t}){const a=w.value;if(!a)return;t&&a.classList.add("transition-disabled");const s=H();s&&ee.value&&(ee.value.style.width=`${s.offsetWidth}px`,ee.value.style.height=`${s.offsetHeight}px`,ee.value.style.transform=`translateX(${s.offsetLeft-Xt(getComputedStyle(a).paddingLeft)}px)`,t&&ee.value.offsetWidth),t&&a.classList.remove("transition-disabled")}fe([V],()=>{e.type==="segment"&&me(()=>{ue({transitionDisabled:!1})})}),He(()=>{e.type==="segment"&&ue({transitionDisabled:!0})});let ze=0;function Je(t){var a;if(t.contentRect.width===0&&t.contentRect.height===0||ze===t.contentRect.width)return;ze=t.contentRect.width;const{type:s}=e;if((s==="line"||s==="bar")&&Re(),s!=="segment"){const{placement:f}=e;ve((f==="top"||f==="bottom"?(a=x.value)===null||a===void 0?void 0:a.$el:_.value)||null)}}const Qe=Se(Je,64);fe([()=>e.justifyContent,()=>e.size],()=>{me(()=>{const{type:t}=e;(t==="line"||t==="bar")&&Re()})});const te=$(!1);function Ze(t){var a;const{target:s,contentRect:{width:f,height:R}}=t,U=s.parentElement.parentElement.offsetWidth,X=s.parentElement.parentElement.offsetHeight,{placement:re}=e;if(!te.value)re==="top"||re==="bottom"?U<f&&(te.value=!0):X<R&&(te.value=!0);else{const{value:oe}=A;if(!oe)return;re==="top"||re==="bottom"?U-f>oe.$el.offsetWidth&&(te.value=!1):X-R>oe.$el.offsetHeight&&(te.value=!1)}ve(((a=x.value)===null||a===void 0?void 0:a.$el)||null)}const et=Se(Ze,64);function tt(){const{onAdd:t}=e;t&&t(),me(()=>{const a=H(),{value:s}=x;!a||!s||s.scrollTo({left:a.offsetLeft,top:0,behavior:"smooth"})})}function ve(t){if(!t)return;const{placement:a}=e;if(a==="top"||a==="bottom"){const{scrollLeft:s,scrollWidth:f,offsetWidth:R}=t;P.value=s<=0,M.value=s+R>=f}else{const{scrollTop:s,scrollHeight:f,offsetHeight:R}=t;P.value=s<=0,M.value=s+R>=f}}const at=Se(t=>{ve(t.target)},64);It(Pe,{triggerRef:K(e,"trigger"),tabStyleRef:K(e,"tabStyle"),tabClassRef:K(e,"tabClass"),addTabStyleRef:K(e,"addTabStyle"),addTabClassRef:K(e,"addTabClass"),paneClassRef:K(e,"paneClass"),paneStyleRef:K(e,"paneStyle"),mergedClsPrefixRef:C,typeRef:K(e,"type"),closableRef:K(e,"closable"),valueRef:V,tabChangeIdRef:S,onBeforeLeaveRef:K(e,"onBeforeLeave"),activateTab:Ke,handleClose:Ye,handleAdd:tt}),Ft(()=>{g(),T()}),Nt(()=>{const{value:t}=L;if(!t)return;const{value:a}=C,s=`${a}-tabs-nav-scroll-wrapper--shadow-start`,f=`${a}-tabs-nav-scroll-wrapper--shadow-end`;P.value?t.classList.remove(s):t.classList.add(s),M.value?t.classList.remove(f):t.classList.add(f)});const rt={syncBarPosition:()=>{g()}},ot=()=>{ue({transitionDisabled:!0})},$e=pe(()=>{const{value:t}=O,{type:a}=e,s={card:"Card",bar:"Bar",line:"Line",segment:"Segment"}[a],f=`${t}${s}`,{self:{barColor:R,closeIconColor:U,closeIconColorHover:X,closeIconColorPressed:re,tabColor:oe,tabBorderColor:nt,paneTextColor:it,tabFontWeight:st,tabBorderRadius:lt,tabFontWeightActive:dt,colorSegment:ct,fontWeightStrong:bt,tabColorSegment:pt,closeSize:ft,closeIconSize:ut,closeColorHover:vt,closeColorPressed:gt,closeBorderRadius:ht,[q("panePadding",t)]:se,[q("tabPadding",f)]:mt,[q("tabPaddingVertical",f)]:xt,[q("tabGap",f)]:yt,[q("tabGap",`${f}Vertical`)]:St,[q("tabTextColor",a)]:Ct,[q("tabTextColorActive",a)]:wt,[q("tabTextColorHover",a)]:_t,[q("tabTextColorDisabled",a)]:Pt,[q("tabFontSize",t)]:Tt},common:{cubicBezierEaseInOut:kt}}=k.value;return{"--n-bezier":kt,"--n-color-segment":ct,"--n-bar-color":R,"--n-tab-font-size":Tt,"--n-tab-text-color":Ct,"--n-tab-text-color-active":wt,"--n-tab-text-color-disabled":Pt,"--n-tab-text-color-hover":_t,"--n-pane-text-color":it,"--n-tab-border-color":nt,"--n-tab-border-radius":lt,"--n-close-size":ft,"--n-close-icon-size":ut,"--n-close-color-hover":vt,"--n-close-color-pressed":gt,"--n-close-border-radius":ht,"--n-close-icon-color":U,"--n-close-icon-color-hover":X,"--n-close-icon-color-pressed":re,"--n-tab-color":oe,"--n-tab-font-weight":st,"--n-tab-font-weight-active":dt,"--n-tab-padding":mt,"--n-tab-padding-vertical":xt,"--n-tab-gap":yt,"--n-tab-gap-vertical":St,"--n-pane-padding-left":le(se,"left"),"--n-pane-padding-right":le(se,"right"),"--n-pane-padding-top":le(se,"top"),"--n-pane-padding-bottom":le(se,"bottom"),"--n-font-weight-strong":bt,"--n-tab-color-segment":pt}}),ae=h?Ut("tabs",pe(()=>`${O.value[0]}${e.type[0]}`),$e,e):void 0;return Object.assign({mergedClsPrefix:C,mergedValue:V,renderedNames:new Set,segmentCapsuleElRef:ee,tabsPaneWrapperRef:N,tabsElRef:w,barElRef:m,addTabInstRef:A,xScrollInstRef:x,scrollWrapperElRef:L,addTabFixed:te,tabWrapperStyle:G,handleNavResize:Qe,mergedSize:O,handleScroll:at,handleTabsResize:et,cssVars:h?void 0:$e,themeClass:ae==null?void 0:ae.themeClass,animationDirection:ke,renderNameListRef:Te,yScrollElRef:_,handleSegmentResize:ot,onAnimationBeforeLeave:Ne,onAnimationEnter:Ue,onAnimationAfterEnter:Xe,onRender:ae==null?void 0:ae.onRender},rt)},render(){const{mergedClsPrefix:e,type:n,placement:b,addTabFixed:l,addable:d,mergedSize:y,renderNameListRef:C,onRender:h,paneWrapperClass:k,paneWrapperStyle:w,$slots:{default:m,prefix:L,suffix:A}}=this;h==null||h();const x=m?ge(m()).filter(S=>S.type.__TAB_PANE__===!0):[],_=m?ge(m()).filter(S=>S.type.__TAB__===!0):[],P=!_.length,M=n==="card",O=n==="segment",j=!M&&!O&&this.justifyContent;C.value=[];const Y=()=>{const S=v("div",{style:this.tabWrapperStyle,class:`${e}-tabs-wrapper`},j?null:v("div",{class:`${e}-tabs-scroll-padding`,style:b==="top"||b==="bottom"?{width:`${this.tabsPadding}px`}:{height:`${this.tabsPadding}px`}}),P?x.map((G,H)=>(C.value.push(G.props.name),Ce(v(we,Object.assign({},G.props,{internalCreatedByPane:!0,internalLeftPadded:H!==0&&(!j||j==="center"||j==="start"||j==="end")}),G.children?{default:G.children.tab}:void 0)))):_.map((G,H)=>(C.value.push(G.props.name),Ce(H!==0&&!j?Ge(G):G))),!l&&d&&M?je(d,(P?x.length:_.length)!==0):null,j?null:v("div",{class:`${e}-tabs-scroll-padding`,style:{width:`${this.tabsPadding}px`}}));return v("div",{ref:"tabsElRef",class:`${e}-tabs-nav-scroll-content`},M&&d?v(he,{onResize:this.handleTabsResize},{default:()=>S}):S,M?v("div",{class:`${e}-tabs-pad`}):null,M?null:v("div",{ref:"barElRef",class:`${e}-tabs-bar`}))},V=O?"top":b;return v("div",{class:[`${e}-tabs`,this.themeClass,`${e}-tabs--${n}-type`,`${e}-tabs--${y}-size`,j&&`${e}-tabs--flex`,`${e}-tabs--${V}`],style:this.cssVars},v("div",{class:[`${e}-tabs-nav--${n}-type`,`${e}-tabs-nav--${V}`,`${e}-tabs-nav`]},Be(L,S=>S&&v("div",{class:`${e}-tabs-nav__prefix`},S)),O?v(he,{onResize:this.handleSegmentResize},{default:()=>v("div",{class:`${e}-tabs-rail`,ref:"tabsElRef"},v("div",{class:`${e}-tabs-capsule`,ref:"segmentCapsuleElRef"},v("div",{class:`${e}-tabs-wrapper`},v("div",{class:`${e}-tabs-tab`}))),P?x.map((S,G)=>(C.value.push(S.props.name),v(we,Object.assign({},S.props,{internalCreatedByPane:!0,internalLeftPadded:G!==0}),S.children?{default:S.children.tab}:void 0))):_.map((S,G)=>(C.value.push(S.props.name),G===0?S:Ge(S))))}):v(he,{onResize:this.handleNavResize},{default:()=>v("div",{class:`${e}-tabs-nav-scroll-wrapper`,ref:"scrollWrapperElRef"},["top","bottom"].includes(V)?v(ba,{ref:"xScrollInstRef",onScroll:this.handleScroll},{default:Y}):v("div",{class:`${e}-tabs-nav-y-scroll`,onScroll:this.handleScroll,ref:"yScrollElRef"},Y()))}),l&&d&&M?je(d,!0):null,Be(A,S=>S&&v("div",{class:`${e}-tabs-nav__suffix`},S))),P&&(this.animated&&(V==="top"||V==="bottom")?v("div",{ref:"tabsPaneWrapperRef",style:w,class:[`${e}-tabs-pane-wrapper`,k]},Ee(x,this.mergedValue,this.renderedNames,this.onAnimationBeforeLeave,this.onAnimationEnter,this.onAnimationAfterEnter,this.animationDirection)):Ee(x,this.mergedValue,this.renderedNames)))}});function Ee(e,n,b,l,d,y,C){const h=[];return e.forEach(k=>{const{name:w,displayDirective:m,"display-directive":L}=k.props,A=_=>m===_||L===_,x=n===w;if(k.key!==void 0&&(k.key=w),x||A("show")||A("show:lazy")&&b.has(w)){b.has(w)||b.add(w);const _=!A("if");h.push(_?Kt(k,[[qt,x]]):k)}}),C?v(Yt,{name:`${C}-transition`,onBeforeLeave:l,onEnter:d,onAfterEnter:y},{default:()=>h}):h}function je(e,n){return v(we,{ref:"addTabInstRef",key:"__addable",name:"__addable",internalCreatedByPane:!0,internalAddable:!0,internalLeftPadded:n,disabled:typeof e=="object"&&e.disabled})}function Ge(e){const n=Jt(e);return n.props?n.props.internalLeftPadded=!0:n.props={internalLeftPadded:!0},n}function Ce(e){return Array.isArray(e.dynamicProps)?e.dynamicProps.includes("internalLeftPadded")||e.dynamicProps.push("internalLeftPadded"):e.dynamicProps=["internalLeftPadded"],e}const Fe=1357375695;function Ve(e){return _e(e)||Qt(e.kind)}async function ya(e){const n=await Zt(e,Fe,320);return n?{ok:!0,body:n.slice(0,180)}:{ok:!1,error:"无播放地址"}}async function Sa(e){const n=await ea(e,Fe,320);if(n)return{ok:!0,body:n.slice(0,180)};const b=await ta(e,{name:"海阔天空",artist:"Beyond"},320);return b?{ok:!0,body:b.slice(0,180)}:{ok:!1,error:"无播放地址"}}async function Ca(e){const n=_e(e),b=performance.now();try{const l=e.kind==="gdapi"?await ya(n):await Sa(n),d=Math.round(performance.now()-b);return l.ok?{ok:!0,latencyMs:d,body:l.body}:{ok:!1,latencyMs:d,error:l.error,body:l.body}}catch(l){let d=l instanceof Error?l.message:String(l);return/failed to fetch|networkerror|load failed/i.test(d)&&(d="Failed to fetch（多半未放行或源已挂）"),{ok:!1,latencyMs:Math.round(performance.now()-b),error:d}}}const wa={class:"monitor-page"},_a={class:"monitor-page__header"},Pa=["disabled"],Ta={class:"monitor-cards"},ka={class:"monitor-card__top"},Ra={class:"monitor-card__title"},za={class:"monitor-card__url"},$a={class:"monitor-card__badge"},La={class:"monitor-card__meta"},Ba={key:0,class:"monitor-card__error"},Aa={key:1,class:"monitor-card__body"},Wa=["disabled","onClick"],Ma={class:"monitor-cards"},Ea={class:"monitor-card__top"},ja={class:"monitor-card__title"},Ga={class:"monitor-card__url"},Va={class:"monitor-card__badge"},Oa={class:"monitor-field"},Da={class:"monitor-card__meta"},Ha={key:0},Ia={key:0,class:"monitor-card__error"},Fa={key:1,class:"monitor-card__body"},Na={class:"monitor-card__actions"},Ua=["disabled","onClick"],Xa=["disabled","onClick"],Ka=["onClick"],qa=["disabled","onClick"],Ya=["disabled","onClick"],Ja=["onClick"],Qa=["disabled"],F="play-sources",Za=ie({__name:"MonitorPage",setup(e){const n=[{id:"onij.fun",label:"onij.fun",cards:[{id:"ping",title:"Ping",url:"https://onij.fun/go/ping"}]},{id:"172.26.250.1",label:"172.26.250.1",cards:[{id:"uvr",title:"UVR",url:"http://172.26.250.1:5555/health"},{id:"whisper",title:"Whisper",url:"http://172.26.250.1:5555/api/v1/whisper/health"}]}],b=[{label:"GD 公开 API",value:"gdapi"},{label:"GD 播放器",value:"gdstudio"}],l=aa(),d=$(n[0].id),y=$(!1),C=$(""),h=ra({}),k=$(""),w=$("gdapi"),m=$("");function L(u,i){return`${u}::${i}`}function A(u,i){var r;return((r=h[L(u,i)])==null?void 0:r.status)??"idle"}function x(u,i){var r;return((r=h[L(u,i)])==null?void 0:r.error)??""}function _(u,i){var r;return((r=h[L(u,i)])==null?void 0:r.body)??""}function P(u,i){var g;const r=(g=h[L(u,i)])==null?void 0:g.latencyMs;return typeof r=="number"?`${r} ms`:"—"}function M(u,i){var g;const r=(g=h[L(u,i)])==null?void 0:g.checkedAt;return r?new Date(r).toLocaleTimeString():"—"}function O(u){switch(u){case"loading":return"检测中";case"ok":return"正常";case"fail":return"异常";default:return"未检测"}}async function j(u,i){const r=L(u,i.id);h[r]={status:"loading"};const g=performance.now();try{const T=await fetch(i.url,{method:"GET",cache:"no-store",signal:AbortSignal.timeout(8e3)}),N=Math.round(performance.now()-g),Z=(await T.text()).slice(0,500);if(!T.ok){h[r]={status:"fail",latencyMs:N,checkedAt:Date.now(),error:`HTTP ${T.status}`,body:Z||void 0};return}h[r]={status:"ok",latencyMs:N,checkedAt:Date.now(),body:Z||void 0}}catch(T){let N=T instanceof Error?T.message:String(T);/failed to fetch|networkerror|load failed/i.test(N)&&(N="Failed to fetch（多半未放行）"),h[r]={status:"fail",latencyMs:Math.round(performance.now()-g),checkedAt:Date.now(),error:N}}}async function Y(u){const i=L(F,u.id);h[i]={status:"loading"};const r=await Ca(u);h[i]={status:r.ok?"ok":"fail",latencyMs:r.latencyMs,checkedAt:Date.now(),error:r.error,body:r.body}}function V(u){return u==="gdstudio"?"GD 播放器":"GD 公开 API"}async function S(u){if(u.kind!=="gdstudio")return;const i=L(F,u.id);C.value=u.id,h[i]={status:"loading"};const r=performance.now();try{const g=await l.syncGdstudioSign(_e(u));h[i]={status:"ok",latencyMs:Math.round(performance.now()-r),checkedAt:Date.now(),body:`已同步签名 ${g}，正在检测取链…`},await Y(u)}catch(g){const T=g instanceof Error?g.message:String(g);h[i]={status:"fail",latencyMs:Math.round(performance.now()-r),checkedAt:Date.now(),error:T||"同步签名失败"}}finally{C.value=""}}function G(){const u=m.value.trim();u&&(l.addSource({title:k.value,kind:w.value,origin:u}),k.value="",m.value="")}async function H(){y.value=!0;try{if(d.value===F){await Promise.all(l.sources.map(i=>Y(i)));return}const u=n.find(i=>i.id===d.value);if(!u)return;await Promise.all(u.cards.map(i=>j(u.id,i)))}finally{y.value=!1}}return fe(d,()=>{H()}),He(async()=>{await l.loadFromServer(),H()}),(u,i)=>(D(),oa(sa,null,{default:ce(()=>[c("div",wa,[c("div",_a,[i[5]||(i[5]=c("h1",{class:"browse-section-title border-l-4 border-blue-500 pl-2"},"监控",-1)),c("button",{type:"button",class:"monitor-btn",disabled:y.value,onClick:H},B(y.value?"检测中…":"刷新当前标签"),9,Pa)]),Q(W(xa),{value:d.value,"onUpdate:value":i[4]||(i[4]=r=>d.value=r),type:"line",animated:"",class:"monitor-tabs"},{default:ce(()=>[(D(),I(be,null,xe(n,r=>Q(W(Me),{key:r.id,name:r.id,tab:r.label},{default:ce(()=>[c("div",Ta,[(D(!0),I(be,null,xe(r.cards,g=>(D(),I("article",{key:g.id,class:We(["monitor-card",`monitor-card--${A(r.id,g.id)}`])},[c("div",ka,[c("div",null,[c("h2",Ra,B(g.title),1),c("p",za,B(g.url),1)]),c("span",$a,B(O(A(r.id,g.id))),1)]),c("dl",La,[c("div",null,[i[6]||(i[6]=c("dt",null,"延迟",-1)),c("dd",null,B(P(r.id,g.id)),1)]),c("div",null,[i[7]||(i[7]=c("dt",null,"最近检测",-1)),c("dd",null,B(M(r.id,g.id)),1)])]),x(r.id,g.id)?(D(),I("p",Ba,B(x(r.id,g.id)),1)):_(r.id,g.id)?(D(),I("pre",Aa,B(_(r.id,g.id)),1)):ne("",!0),c("button",{type:"button",class:"monitor-btn monitor-btn--ghost",disabled:y.value,onClick:T=>j(r.id,g)}," 单独检测 ",8,Wa)],2))),128))])]),_:2},1032,["name","tab"])),64)),Q(W(Me),{name:F,tab:"音乐源"},{default:ce(()=>[i[13]||(i[13]=c("p",{class:"monitor-hint"}," 取链按卡片从上到下尝试。GD 播放器解码会轮换：先「单独检测」，失败再点「同步签名」，会从源站拉最新 player.js / crc32 写到服务器，全站立刻用新解码，不用改代码发版。 ",-1)),c("div",Ma,[(D(!0),I(be,null,xe(W(l).sources,(r,g)=>(D(),I("article",{key:r.id,class:We(["monitor-card",[`monitor-card--${A(F,r.id)}`,{"monitor-card--off":!r.enabled}]])},[c("div",Ea,[c("div",null,[c("h2",ja,B(r.title),1),c("p",Ga,B(V(r.kind))+" · "+B(W(Ve)(r)),1)]),c("span",Va,B(O(A(F,r.id))),1)]),c("label",Oa,[i[8]||(i[8]=c("span",null,"源地址",-1)),Q(W(ye),{value:r.origin,placeholder:W(Ve)(r),size:"small","onUpdate:value":T=>W(l).patchSource(r.id,{origin:T}),onBlur:i[0]||(i[0]=T=>W(l).flushSave())},null,8,["value","placeholder","onUpdate:value"])]),c("dl",Da,[c("div",null,[i[9]||(i[9]=c("dt",null,"延迟",-1)),c("dd",null,B(P(F,r.id)),1)]),c("div",null,[i[10]||(i[10]=c("dt",null,"最近检测",-1)),c("dd",null,B(M(F,r.id)),1)]),r.kind==="gdstudio"?(D(),I("div",Ha,[i[11]||(i[11]=c("dt",null,"签名版本",-1)),c("dd",null,B(W(l).signVersion||"未同步"),1)])):ne("",!0)]),x(F,r.id)?(D(),I("p",Ia,B(x(F,r.id)),1)):_(F,r.id)?(D(),I("pre",Fa,B(_(F,r.id)),1)):ne("",!0),c("div",Na,[c("button",{type:"button",class:"monitor-btn monitor-btn--ghost",disabled:y.value,onClick:T=>Y(r)}," 单独检测 ",8,Ua),r.kind==="gdstudio"?(D(),I("button",{key:0,type:"button",class:"monitor-btn monitor-btn--ghost",disabled:y.value||C.value===r.id,onClick:T=>S(r)},B(C.value===r.id?"同步中…":"同步签名"),9,Xa)):ne("",!0),c("button",{type:"button",class:"monitor-btn monitor-btn--ghost",onClick:T=>W(l).patchSource(r.id,{enabled:!r.enabled})},B(r.enabled?"停用":"启用"),9,Ka),c("button",{type:"button",class:"monitor-btn monitor-btn--ghost",disabled:g===0,onClick:T=>W(l).moveSource(r.id,-1)}," 上移 ",8,qa),c("button",{type:"button",class:"monitor-btn monitor-btn--ghost",disabled:g===W(l).sources.length-1,onClick:T=>W(l).moveSource(r.id,1)}," 下移 ",8,Ya),r.builtin?ne("",!0):(D(),I("button",{key:1,type:"button",class:"monitor-btn monitor-btn--ghost",onClick:T=>W(l).removeSource(r.id)}," 删除 ",8,Ja))])],2))),128))]),c("form",{class:"monitor-add",onSubmit:na(G,["prevent"])},[i[12]||(i[12]=c("h3",{class:"monitor-add__title"},"添加源",-1)),Q(W(ye),{value:k.value,"onUpdate:value":i[1]||(i[1]=r=>k.value=r),size:"small",placeholder:"名称"},null,8,["value"]),Q(W(da),{value:w.value,"onUpdate:value":i[2]||(i[2]=r=>w.value=r),size:"small",options:b},null,8,["value"]),Q(W(ye),{value:m.value,"onUpdate:value":i[3]||(i[3]=r=>m.value=r),size:"small",placeholder:"https://music-api.gdstudio.xyz"},null,8,["value"]),c("button",{type:"submit",class:"monitor-btn",disabled:!m.value.trim()},"添加",8,Qa)],32)]),_:1})]),_:1},8,["value"])])]),_:1}))}}),rr=ia(Za,[["__scopeId","data-v-a3583e79"]]);export{rr as default};
