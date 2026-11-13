import{d as re,aa as g,r as L,bK as Le,bL as kt,bM as zt,bN as $t,bO as Lt,ab as Bt,av as At,a_ as Mt,a5 as Ge,aN as Wt,b3 as Et,F as de,aF as jt,ar as Ot,c as ce,az as Vt,ad as o,ah as c,ac as $,ai as M,aD as Gt,bl as ve,aj as Be,bP as ge,ak as Ht,al as He,bx as Ae,aG as Dt,w as be,o as De,aC as Ft,aw as N,bQ as It,aH as Nt,an as U,ao as ie,ap as Ut,bo as Xt,A as Kt,bc as Yt,bR as qt,ba as Jt,ax as he,as as se,bS as _e,bT as Qt,bU as Zt,bV as ea,bW as ta,a as aa,e as ra,f as le,i as b,n as q,t as W,g as X,j as me,y as j,m as oa,l as F,I as Me,h as xe,_ as na}from"./index-DmwaRmJO.js";import{P as ia}from"./PageContent-Pd-l6AUA.js";import{A as sa,d as ye,f as la}from"./Select-BBs4N5Km.js";const da=Le(".v-x-scroll",{overflow:"auto",scrollbarWidth:"none"},[Le("&::-webkit-scrollbar",{width:0,height:0})]),ca=re({name:"XScroll",props:{disabled:Boolean,onScroll:Function},setup(){const e=L(null);function n(f){!(f.currentTarget.offsetWidth<f.currentTarget.scrollWidth)||f.deltaY===0||(f.currentTarget.scrollLeft+=f.deltaY+f.deltaX,f.preventDefault())}const l=kt();return da.mount({id:"vueuc/x-scroll",head:!0,anchorMetaName:zt,ssr:l}),Object.assign({selfRef:e,handleWheel:n},{scrollTo(...f){var m;(m=e.value)===null||m===void 0||m.scrollTo(...f)}})},render(){return g("div",{ref:"selfRef",onScroll:this.onScroll,onWheel:this.disabled?void 0:this.handleWheel,class:"v-x-scroll"},this.$slots)}});var ba="Expected a function";function Se(e,n,l){var p=!0,f=!0;if(typeof e!="function")throw new TypeError(ba);return Lt(l)&&(p="leading"in l?!!l.leading:p,f="trailing"in l?!!l.trailing:f),$t(e,n,{leading:p,maxWait:n,trailing:f})}const pa={tabFontSizeSmall:"14px",tabFontSizeMedium:"14px",tabFontSizeLarge:"16px",tabGapSmallLine:"36px",tabGapMediumLine:"36px",tabGapLargeLine:"36px",tabGapSmallLineVertical:"8px",tabGapMediumLineVertical:"8px",tabGapLargeLineVertical:"8px",tabPaddingSmallLine:"6px 0",tabPaddingMediumLine:"10px 0",tabPaddingLargeLine:"14px 0",tabPaddingVerticalSmallLine:"6px 12px",tabPaddingVerticalMediumLine:"8px 16px",tabPaddingVerticalLargeLine:"10px 20px",tabGapSmallBar:"36px",tabGapMediumBar:"36px",tabGapLargeBar:"36px",tabGapSmallBarVertical:"8px",tabGapMediumBarVertical:"8px",tabGapLargeBarVertical:"8px",tabPaddingSmallBar:"4px 0",tabPaddingMediumBar:"6px 0",tabPaddingLargeBar:"10px 0",tabPaddingVerticalSmallBar:"6px 12px",tabPaddingVerticalMediumBar:"8px 16px",tabPaddingVerticalLargeBar:"10px 20px",tabGapSmallCard:"4px",tabGapMediumCard:"4px",tabGapLargeCard:"4px",tabGapSmallCardVertical:"4px",tabGapMediumCardVertical:"4px",tabGapLargeCardVertical:"4px",tabPaddingSmallCard:"8px 16px",tabPaddingMediumCard:"10px 20px",tabPaddingLargeCard:"12px 24px",tabPaddingSmallSegment:"4px 0",tabPaddingMediumSegment:"6px 0",tabPaddingLargeSegment:"8px 0",tabPaddingVerticalLargeSegment:"0 8px",tabPaddingVerticalSmallCard:"8px 12px",tabPaddingVerticalMediumCard:"10px 16px",tabPaddingVerticalLargeCard:"12px 20px",tabPaddingVerticalSmallSegment:"0 4px",tabPaddingVerticalMediumSegment:"0 6px",tabGapSmallSegment:"0",tabGapMediumSegment:"0",tabGapLargeSegment:"0",tabGapSmallSegmentVertical:"0",tabGapMediumSegmentVertical:"0",tabGapLargeSegmentVertical:"0",panePaddingSmall:"8px 0 0 0",panePaddingMedium:"12px 0 0 0",panePaddingLarge:"16px 0 0 0",closeSize:"18px",closeIconSize:"14px"};function fa(e){const{textColor2:n,primaryColor:l,textColorDisabled:p,closeIconColor:f,closeIconColorHover:m,closeIconColorPressed:u,closeColorHover:C,closeColorPressed:_,tabColor:x,baseColor:h,dividerColor:R,fontWeight:T,textColor1:y,borderRadius:k,fontSize:P,fontWeightStrong:A}=e;return Object.assign(Object.assign({},pa),{colorSegment:x,tabFontSizeCard:P,tabTextColorLine:y,tabTextColorActiveLine:l,tabTextColorHoverLine:l,tabTextColorDisabledLine:p,tabTextColorSegment:y,tabTextColorActiveSegment:n,tabTextColorHoverSegment:n,tabTextColorDisabledSegment:p,tabTextColorBar:y,tabTextColorActiveBar:l,tabTextColorHoverBar:l,tabTextColorDisabledBar:p,tabTextColorCard:y,tabTextColorHoverCard:y,tabTextColorActiveCard:l,tabTextColorDisabledCard:p,barColor:l,closeIconColor:f,closeIconColorHover:m,closeIconColorPressed:u,closeColorHover:C,closeColorPressed:_,closeBorderRadius:k,tabColor:x,tabColorSegment:h,tabBorderColor:R,tabFontWeightActive:T,tabFontWeight:T,tabBorderRadius:k,paneTextColor:n,fontWeightStrong:A})}const ua={common:Bt,self:fa},Pe=At("n-tabs"),Fe={tab:[String,Number,Object,Function],name:{type:[String,Number],required:!0},disabled:Boolean,displayDirective:{type:String,default:"if"},closable:{type:Boolean,default:void 0},tabProps:Object,label:[String,Number,Object,Function]},We=re({__TAB_PANE__:!0,name:"TabPane",alias:["TabPanel"],props:Fe,slots:Object,setup(e){const n=Ge(Pe,null);return n||Mt("tab-pane","`n-tab-pane` must be placed inside `n-tabs`."),{style:n.paneStyleRef,class:n.paneClassRef,mergedClsPrefix:n.mergedClsPrefixRef}},render(){return g("div",{class:[`${this.mergedClsPrefix}-tab-pane`,this.class],style:this.style},this.$slots)}}),va=Object.assign({internalLeftPadded:Boolean,internalAddable:Boolean,internalCreatedByPane:Boolean},Vt(Fe,["displayDirective"])),we=re({__TAB__:!0,inheritAttrs:!1,name:"Tab",props:va,setup(e){const{mergedClsPrefixRef:n,valueRef:l,typeRef:p,closableRef:f,tabStyleRef:m,addTabStyleRef:u,tabClassRef:C,addTabClassRef:_,tabChangeIdRef:x,onBeforeLeaveRef:h,triggerRef:R,handleAdd:T,activateTab:y,handleClose:k}=Ge(Pe);return{trigger:R,mergedClosable:ce(()=>{if(e.internalAddable)return!1;const{closable:P}=e;return P===void 0?f.value:P}),style:m,addStyle:u,tabClass:C,addTabClass:_,clsPrefix:n,value:l,type:p,handleClose(P){P.stopPropagation(),!e.disabled&&k(e.name)},activateTab(){if(e.disabled)return;if(e.internalAddable){T();return}const{name:P}=e,A=++x.id;if(P!==l.value){const{value:V}=h;V?Promise.resolve(V(e.name,l.value)).then(E=>{E&&x.id===A&&y(P)}):y(P)}}}},render(){const{internalAddable:e,clsPrefix:n,name:l,disabled:p,label:f,tab:m,value:u,mergedClosable:C,trigger:_,$slots:{default:x}}=this,h=f??m;return g("div",{class:`${n}-tabs-tab-wrapper`},this.internalLeftPadded?g("div",{class:`${n}-tabs-tab-pad`}):null,g("div",Object.assign({key:l,"data-name":l,"data-disabled":p?!0:void 0},Wt({class:[`${n}-tabs-tab`,u===l&&`${n}-tabs-tab--active`,p&&`${n}-tabs-tab--disabled`,C&&`${n}-tabs-tab--closable`,e&&`${n}-tabs-tab--addable`,e?this.addTabClass:this.tabClass],onClick:_==="click"?this.activateTab:void 0,onMouseenter:_==="hover"?this.activateTab:void 0,style:e?this.addStyle:this.style},this.internalCreatedByPane?this.tabProps||{}:this.$attrs)),g("span",{class:`${n}-tabs-tab__label`},e?g(de,null,g("div",{class:`${n}-tabs-tab__height-placeholder`}," "),g(jt,{clsPrefix:n},{default:()=>g(sa,null)})):x?x():typeof h=="object"?h:Et(h??l)),C&&this.type==="card"?g(Ot,{clsPrefix:n,class:`${n}-tabs-tab__close`,onClick:this.handleClose,disabled:p}):null))}}),ga=o("tabs",`
 box-sizing: border-box;
 width: 100%;
 display: flex;
 flex-direction: column;
 transition:
 background-color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
`,[c("segment-type",[o("tabs-rail",[$("&.transition-disabled",[o("tabs-capsule",`
 transition: none;
 `)])])]),c("top",[o("tab-pane",`
 padding: var(--n-pane-padding-top) var(--n-pane-padding-right) var(--n-pane-padding-bottom) var(--n-pane-padding-left);
 `)]),c("left",[o("tab-pane",`
 padding: var(--n-pane-padding-right) var(--n-pane-padding-bottom) var(--n-pane-padding-left) var(--n-pane-padding-top);
 `)]),c("left, right",`
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
 `)]),c("right",`
 flex-direction: row-reverse;
 `,[o("tab-pane",`
 padding: var(--n-pane-padding-left) var(--n-pane-padding-top) var(--n-pane-padding-right) var(--n-pane-padding-bottom);
 `),o("tabs-bar",`
 left: 0;
 `)]),c("bottom",`
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
 `,[c("active",`
 font-weight: var(--n-font-weight-strong);
 color: var(--n-tab-text-color-active);
 `),$("&:hover",`
 color: var(--n-tab-text-color-hover);
 `)])])]),c("flex",[o("tabs-nav",`
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
 `,[M("prefix, suffix",`
 display: flex;
 align-items: center;
 `),M("prefix","padding-right: 16px;"),M("suffix","padding-left: 16px;")]),c("top, bottom",[o("tabs-nav-scroll-wrapper",[$("&::before",`
 top: 0;
 bottom: 0;
 left: 0;
 width: 20px;
 `),$("&::after",`
 top: 0;
 bottom: 0;
 right: 0;
 width: 20px;
 `),c("shadow-start",[$("&::before",`
 box-shadow: inset 10px 0 8px -8px rgba(0, 0, 0, .12);
 `)]),c("shadow-end",[$("&::after",`
 box-shadow: inset -10px 0 8px -8px rgba(0, 0, 0, .12);
 `)])])]),c("left, right",[o("tabs-nav-scroll-content",`
 flex-direction: column;
 `),o("tabs-nav-scroll-wrapper",[$("&::before",`
 top: 0;
 left: 0;
 right: 0;
 height: 20px;
 `),$("&::after",`
 bottom: 0;
 left: 0;
 right: 0;
 height: 20px;
 `),c("shadow-start",[$("&::before",`
 box-shadow: inset 0 10px 8px -8px rgba(0, 0, 0, .12);
 `)]),c("shadow-end",[$("&::after",`
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
 `,[$("&::-webkit-scrollbar, &::-webkit-scrollbar-track-piece, &::-webkit-scrollbar-thumb",`
 width: 0;
 height: 0;
 display: none;
 `)]),$("&::before, &::after",`
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
 `,[c("disabled",{cursor:"not-allowed"}),M("close",`
 margin-left: 6px;
 transition:
 background-color .3s var(--n-bezier),
 color .3s var(--n-bezier);
 `),M("label",`
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
 `,[$("&.transition-disabled",`
 transition: none;
 `),c("disabled",`
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
 `,[$("&.next-transition-leave-active, &.prev-transition-leave-active, &.next-transition-enter-active, &.prev-transition-enter-active",`
 transition:
 color .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 transform .2s var(--n-bezier),
 opacity .2s var(--n-bezier);
 `),$("&.next-transition-leave-active, &.prev-transition-leave-active",`
 position: absolute;
 `),$("&.next-transition-enter-from, &.prev-transition-leave-to",`
 transform: translateX(32px);
 opacity: 0;
 `),$("&.next-transition-leave-to, &.prev-transition-enter-from",`
 transform: translateX(-32px);
 opacity: 0;
 `),$("&.next-transition-leave-from, &.next-transition-enter-to, &.prev-transition-leave-from, &.prev-transition-enter-to",`
 transform: translateX(0);
 opacity: 1;
 `)]),o("tabs-tab-pad",`
 box-sizing: border-box;
 width: var(--n-tab-gap);
 flex-grow: 0;
 flex-shrink: 0;
 `),c("line-type, bar-type",[o("tabs-tab",`
 font-weight: var(--n-tab-font-weight);
 box-sizing: border-box;
 vertical-align: bottom;
 `,[$("&:hover",{color:"var(--n-tab-text-color-hover)"}),c("active",`
 color: var(--n-tab-text-color-active);
 font-weight: var(--n-tab-font-weight-active);
 `),c("disabled",{color:"var(--n-tab-text-color-disabled)"})])]),o("tabs-nav",[c("line-type",[c("top",[M("prefix, suffix",`
 border-bottom: 1px solid var(--n-tab-border-color);
 `),o("tabs-nav-scroll-content",`
 border-bottom: 1px solid var(--n-tab-border-color);
 `),o("tabs-bar",`
 bottom: -1px;
 `)]),c("left",[M("prefix, suffix",`
 border-right: 1px solid var(--n-tab-border-color);
 `),o("tabs-nav-scroll-content",`
 border-right: 1px solid var(--n-tab-border-color);
 `),o("tabs-bar",`
 right: -1px;
 `)]),c("right",[M("prefix, suffix",`
 border-left: 1px solid var(--n-tab-border-color);
 `),o("tabs-nav-scroll-content",`
 border-left: 1px solid var(--n-tab-border-color);
 `),o("tabs-bar",`
 left: -1px;
 `)]),c("bottom",[M("prefix, suffix",`
 border-top: 1px solid var(--n-tab-border-color);
 `),o("tabs-nav-scroll-content",`
 border-top: 1px solid var(--n-tab-border-color);
 `),o("tabs-bar",`
 top: -1px;
 `)]),M("prefix, suffix",`
 transition: border-color .3s var(--n-bezier);
 `),o("tabs-nav-scroll-content",`
 transition: border-color .3s var(--n-bezier);
 `),o("tabs-bar",`
 border-radius: 0;
 `)]),c("card-type",[M("prefix, suffix",`
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
 `,[c("addable",`
 padding-left: 8px;
 padding-right: 8px;
 font-size: 16px;
 justify-content: center;
 `,[M("height-placeholder",`
 width: 0;
 font-size: var(--n-tab-font-size);
 `),Gt("disabled",[$("&:hover",`
 color: var(--n-tab-text-color-hover);
 `)])]),c("closable","padding-right: 8px;"),c("active",`
 background-color: #0000;
 font-weight: var(--n-tab-font-weight-active);
 color: var(--n-tab-text-color-active);
 `),c("disabled","color: var(--n-tab-text-color-disabled);")])]),c("left, right",`
 flex-direction: column; 
 `,[M("prefix, suffix",`
 padding: var(--n-tab-padding-vertical);
 `),o("tabs-wrapper",`
 flex-direction: column;
 `),o("tabs-tab-wrapper",`
 flex-direction: column;
 `,[o("tabs-tab-pad",`
 height: var(--n-tab-gap-vertical);
 width: 100%;
 `)])]),c("top",[c("card-type",[o("tabs-scroll-padding","border-bottom: 1px solid var(--n-tab-border-color);"),M("prefix, suffix",`
 border-bottom: 1px solid var(--n-tab-border-color);
 `),o("tabs-tab",`
 border-top-left-radius: var(--n-tab-border-radius);
 border-top-right-radius: var(--n-tab-border-radius);
 `,[c("active",`
 border-bottom: 1px solid #0000;
 `)]),o("tabs-tab-pad",`
 border-bottom: 1px solid var(--n-tab-border-color);
 `),o("tabs-pad",`
 border-bottom: 1px solid var(--n-tab-border-color);
 `)])]),c("left",[c("card-type",[o("tabs-scroll-padding","border-right: 1px solid var(--n-tab-border-color);"),M("prefix, suffix",`
 border-right: 1px solid var(--n-tab-border-color);
 `),o("tabs-tab",`
 border-top-left-radius: var(--n-tab-border-radius);
 border-bottom-left-radius: var(--n-tab-border-radius);
 `,[c("active",`
 border-right: 1px solid #0000;
 `)]),o("tabs-tab-pad",`
 border-right: 1px solid var(--n-tab-border-color);
 `),o("tabs-pad",`
 border-right: 1px solid var(--n-tab-border-color);
 `)])]),c("right",[c("card-type",[o("tabs-scroll-padding","border-left: 1px solid var(--n-tab-border-color);"),M("prefix, suffix",`
 border-left: 1px solid var(--n-tab-border-color);
 `),o("tabs-tab",`
 border-top-right-radius: var(--n-tab-border-radius);
 border-bottom-right-radius: var(--n-tab-border-radius);
 `,[c("active",`
 border-left: 1px solid #0000;
 `)]),o("tabs-tab-pad",`
 border-left: 1px solid var(--n-tab-border-color);
 `),o("tabs-pad",`
 border-left: 1px solid var(--n-tab-border-color);
 `)])]),c("bottom",[c("card-type",[o("tabs-scroll-padding","border-top: 1px solid var(--n-tab-border-color);"),M("prefix, suffix",`
 border-top: 1px solid var(--n-tab-border-color);
 `),o("tabs-tab",`
 border-bottom-left-radius: var(--n-tab-border-radius);
 border-bottom-right-radius: var(--n-tab-border-radius);
 `,[c("active",`
 border-top: 1px solid #0000;
 `)]),o("tabs-tab-pad",`
 border-top: 1px solid var(--n-tab-border-color);
 `),o("tabs-pad",`
 border-top: 1px solid var(--n-tab-border-color);
 `)])])])]),ha=Object.assign(Object.assign({},He.props),{value:[String,Number],defaultValue:[String,Number],trigger:{type:String,default:"click"},type:{type:String,default:"bar"},closable:Boolean,justifyContent:String,size:{type:String,default:"medium"},placement:{type:String,default:"top"},tabStyle:[String,Object],tabClass:String,addTabStyle:[String,Object],addTabClass:String,barWidth:Number,paneClass:String,paneStyle:[String,Object],paneWrapperClass:String,paneWrapperStyle:[String,Object],addable:[Boolean,Object],tabsPadding:{type:Number,default:0},animated:Boolean,onBeforeLeave:Function,onAdd:Function,"onUpdate:value":[Function,Array],onUpdateValue:[Function,Array],onClose:[Function,Array],labelSize:String,activeName:[String,Number],onActiveNameChange:[Function,Array]}),ma=re({name:"Tabs",props:ha,slots:Object,setup(e,{slots:n}){var l,p,f,m;const{mergedClsPrefixRef:u,inlineThemeDisabled:C}=Ht(e),_=He("Tabs","-tabs",ga,ua,e,u),x=L(null),h=L(null),R=L(null),T=L(null),y=L(null),k=L(null),P=L(!0),A=L(!0),V=Ae(e,["labelSize","size"]),E=Ae(e,["activeName","value"]),Y=L((p=(l=E.value)!==null&&l!==void 0?l:e.defaultValue)!==null&&p!==void 0?p:n.default?(m=(f=ve(n.default())[0])===null||f===void 0?void 0:f.props)===null||m===void 0?void 0:m.name:null),O=Dt(E,Y),S={id:0},d=ce(()=>{if(!(!e.justifyContent||e.type==="card"))return{display:"flex",justifyContent:e.justifyContent}});be(O,()=>{S.id=0,G(),te()});function i(){var t;const{value:a}=O;return a===null?null:(t=x.value)===null||t===void 0?void 0:t.querySelector(`[data-name="${a}"]`)}function r(t){if(e.type==="card")return;const{value:a}=h;if(!a)return;const s=a.style.opacity==="0";if(t){const v=`${u.value}-tabs-bar--disabled`,{barWidth:z,placement:H}=e;if(t.dataset.disabled==="true"?a.classList.add(v):a.classList.remove(v),["top","bottom"].includes(H)){if(B(["top","maxHeight","height"]),typeof z=="number"&&t.offsetWidth>=z){const D=Math.floor((t.offsetWidth-z)/2)+t.offsetLeft;a.style.left=`${D}px`,a.style.maxWidth=`${z}px`}else a.style.left=`${t.offsetLeft}px`,a.style.maxWidth=`${t.offsetWidth}px`;a.style.width="8192px",s&&(a.style.transition="none"),a.offsetWidth,s&&(a.style.transition="",a.style.opacity="1")}else{if(B(["left","maxWidth","width"]),typeof z=="number"&&t.offsetHeight>=z){const D=Math.floor((t.offsetHeight-z)/2)+t.offsetTop;a.style.top=`${D}px`,a.style.maxHeight=`${z}px`}else a.style.top=`${t.offsetTop}px`,a.style.maxHeight=`${t.offsetHeight}px`;a.style.height="8192px",s&&(a.style.transition="none"),a.offsetHeight,s&&(a.style.transition="",a.style.opacity="1")}}}function w(){if(e.type==="card")return;const{value:t}=h;t&&(t.style.opacity="0")}function B(t){const{value:a}=h;if(a)for(const s of t)a.style[s]=""}function G(){if(e.type==="card")return;const t=i();t?r(t):w()}function te(){var t;const a=(t=y.value)===null||t===void 0?void 0:t.$el;if(!a)return;const s=i();if(!s)return;const{scrollLeft:v,offsetWidth:z}=a,{offsetLeft:H,offsetWidth:D}=s;v>H?a.scrollTo({top:0,left:H,behavior:"smooth"}):H+D>v+z&&a.scrollTo({top:0,left:H+D-z,behavior:"smooth"})}const oe=L(null);let pe=0,K=null;function Ne(t){const a=oe.value;if(a){pe=t.getBoundingClientRect().height;const s=`${pe}px`,v=()=>{a.style.height=s,a.style.maxHeight=s};K?(v(),K(),K=null):K=v}}function Ue(t){const a=oe.value;if(a){const s=t.getBoundingClientRect().height,v=()=>{document.body.offsetHeight,a.style.maxHeight=`${s}px`,a.style.height=`${Math.max(pe,s)}px`};K?(K(),K=null,v()):K=v}}function Xe(){const t=oe.value;if(t){t.style.maxHeight="",t.style.height="";const{paneWrapperStyle:a}=e;if(typeof a=="string")t.style.cssText=a;else if(a){const{maxHeight:s,height:v}=a;s!==void 0&&(t.style.maxHeight=s),v!==void 0&&(t.style.height=v)}}}const Te={value:[]},Re=L("next");function Ke(t){const a=O.value;let s="next";for(const v of Te.value){if(v===a)break;if(v===t){s="prev";break}}Re.value=s,Ye(t)}function Ye(t){const{onActiveNameChange:a,onUpdateValue:s,"onUpdate:value":v}=e;a&&se(a,t),s&&se(s,t),v&&se(v,t),Y.value=t}function qe(t){const{onClose:a}=e;a&&se(a,t)}function ke(){const{value:t}=h;if(!t)return;const a="transition-disabled";t.classList.add(a),G(),t.classList.remove(a)}const J=L(null);function fe({transitionDisabled:t}){const a=x.value;if(!a)return;t&&a.classList.add("transition-disabled");const s=i();s&&J.value&&(J.value.style.width=`${s.offsetWidth}px`,J.value.style.height=`${s.offsetHeight}px`,J.value.style.transform=`translateX(${s.offsetLeft-Xt(getComputedStyle(a).paddingLeft)}px)`,t&&J.value.offsetWidth),t&&a.classList.remove("transition-disabled")}be([O],()=>{e.type==="segment"&&he(()=>{fe({transitionDisabled:!1})})}),De(()=>{e.type==="segment"&&fe({transitionDisabled:!0})});let ze=0;function Je(t){var a;if(t.contentRect.width===0&&t.contentRect.height===0||ze===t.contentRect.width)return;ze=t.contentRect.width;const{type:s}=e;if((s==="line"||s==="bar")&&ke(),s!=="segment"){const{placement:v}=e;ue((v==="top"||v==="bottom"?(a=y.value)===null||a===void 0?void 0:a.$el:k.value)||null)}}const Qe=Se(Je,64);be([()=>e.justifyContent,()=>e.size],()=>{he(()=>{const{type:t}=e;(t==="line"||t==="bar")&&ke()})});const Q=L(!1);function Ze(t){var a;const{target:s,contentRect:{width:v,height:z}}=t,H=s.parentElement.parentElement.offsetWidth,D=s.parentElement.parentElement.offsetHeight,{placement:ee}=e;if(!Q.value)ee==="top"||ee==="bottom"?H<v&&(Q.value=!0):D<z&&(Q.value=!0);else{const{value:ae}=T;if(!ae)return;ee==="top"||ee==="bottom"?H-v>ae.$el.offsetWidth&&(Q.value=!1):D-z>ae.$el.offsetHeight&&(Q.value=!1)}ue(((a=y.value)===null||a===void 0?void 0:a.$el)||null)}const et=Se(Ze,64);function tt(){const{onAdd:t}=e;t&&t(),he(()=>{const a=i(),{value:s}=y;!a||!s||s.scrollTo({left:a.offsetLeft,top:0,behavior:"smooth"})})}function ue(t){if(!t)return;const{placement:a}=e;if(a==="top"||a==="bottom"){const{scrollLeft:s,scrollWidth:v,offsetWidth:z}=t;P.value=s<=0,A.value=s+z>=v}else{const{scrollTop:s,scrollHeight:v,offsetHeight:z}=t;P.value=s<=0,A.value=s+z>=v}}const at=Se(t=>{ue(t.target)},64);Ft(Pe,{triggerRef:N(e,"trigger"),tabStyleRef:N(e,"tabStyle"),tabClassRef:N(e,"tabClass"),addTabStyleRef:N(e,"addTabStyle"),addTabClassRef:N(e,"addTabClass"),paneClassRef:N(e,"paneClass"),paneStyleRef:N(e,"paneStyle"),mergedClsPrefixRef:u,typeRef:N(e,"type"),closableRef:N(e,"closable"),valueRef:O,tabChangeIdRef:S,onBeforeLeaveRef:N(e,"onBeforeLeave"),activateTab:Ke,handleClose:qe,handleAdd:tt}),It(()=>{G(),te()}),Nt(()=>{const{value:t}=R;if(!t)return;const{value:a}=u,s=`${a}-tabs-nav-scroll-wrapper--shadow-start`,v=`${a}-tabs-nav-scroll-wrapper--shadow-end`;P.value?t.classList.remove(s):t.classList.add(s),A.value?t.classList.remove(v):t.classList.add(v)});const rt={syncBarPosition:()=>{G()}},ot=()=>{fe({transitionDisabled:!0})},$e=ce(()=>{const{value:t}=V,{type:a}=e,s={card:"Card",bar:"Bar",line:"Line",segment:"Segment"}[a],v=`${t}${s}`,{self:{barColor:z,closeIconColor:H,closeIconColorHover:D,closeIconColorPressed:ee,tabColor:ae,tabBorderColor:nt,paneTextColor:it,tabFontWeight:st,tabBorderRadius:lt,tabFontWeightActive:dt,colorSegment:ct,fontWeightStrong:bt,tabColorSegment:pt,closeSize:ft,closeIconSize:ut,closeColorHover:vt,closeColorPressed:gt,closeBorderRadius:ht,[U("panePadding",t)]:ne,[U("tabPadding",v)]:mt,[U("tabPaddingVertical",v)]:xt,[U("tabGap",v)]:yt,[U("tabGap",`${v}Vertical`)]:St,[U("tabTextColor",a)]:Ct,[U("tabTextColorActive",a)]:wt,[U("tabTextColorHover",a)]:_t,[U("tabTextColorDisabled",a)]:Pt,[U("tabFontSize",t)]:Tt},common:{cubicBezierEaseInOut:Rt}}=_.value;return{"--n-bezier":Rt,"--n-color-segment":ct,"--n-bar-color":z,"--n-tab-font-size":Tt,"--n-tab-text-color":Ct,"--n-tab-text-color-active":wt,"--n-tab-text-color-disabled":Pt,"--n-tab-text-color-hover":_t,"--n-pane-text-color":it,"--n-tab-border-color":nt,"--n-tab-border-radius":lt,"--n-close-size":ft,"--n-close-icon-size":ut,"--n-close-color-hover":vt,"--n-close-color-pressed":gt,"--n-close-border-radius":ht,"--n-close-icon-color":H,"--n-close-icon-color-hover":D,"--n-close-icon-color-pressed":ee,"--n-tab-color":ae,"--n-tab-font-weight":st,"--n-tab-font-weight-active":dt,"--n-tab-padding":mt,"--n-tab-padding-vertical":xt,"--n-tab-gap":yt,"--n-tab-gap-vertical":St,"--n-pane-padding-left":ie(ne,"left"),"--n-pane-padding-right":ie(ne,"right"),"--n-pane-padding-top":ie(ne,"top"),"--n-pane-padding-bottom":ie(ne,"bottom"),"--n-font-weight-strong":bt,"--n-tab-color-segment":pt}}),Z=C?Ut("tabs",ce(()=>`${V.value[0]}${e.type[0]}`),$e,e):void 0;return Object.assign({mergedClsPrefix:u,mergedValue:O,renderedNames:new Set,segmentCapsuleElRef:J,tabsPaneWrapperRef:oe,tabsElRef:x,barElRef:h,addTabInstRef:T,xScrollInstRef:y,scrollWrapperElRef:R,addTabFixed:Q,tabWrapperStyle:d,handleNavResize:Qe,mergedSize:V,handleScroll:at,handleTabsResize:et,cssVars:C?void 0:$e,themeClass:Z==null?void 0:Z.themeClass,animationDirection:Re,renderNameListRef:Te,yScrollElRef:k,handleSegmentResize:ot,onAnimationBeforeLeave:Ne,onAnimationEnter:Ue,onAnimationAfterEnter:Xe,onRender:Z==null?void 0:Z.onRender},rt)},render(){const{mergedClsPrefix:e,type:n,placement:l,addTabFixed:p,addable:f,mergedSize:m,renderNameListRef:u,onRender:C,paneWrapperClass:_,paneWrapperStyle:x,$slots:{default:h,prefix:R,suffix:T}}=this;C==null||C();const y=h?ve(h()).filter(S=>S.type.__TAB_PANE__===!0):[],k=h?ve(h()).filter(S=>S.type.__TAB__===!0):[],P=!k.length,A=n==="card",V=n==="segment",E=!A&&!V&&this.justifyContent;u.value=[];const Y=()=>{const S=g("div",{style:this.tabWrapperStyle,class:`${e}-tabs-wrapper`},E?null:g("div",{class:`${e}-tabs-scroll-padding`,style:l==="top"||l==="bottom"?{width:`${this.tabsPadding}px`}:{height:`${this.tabsPadding}px`}}),P?y.map((d,i)=>(u.value.push(d.props.name),Ce(g(we,Object.assign({},d.props,{internalCreatedByPane:!0,internalLeftPadded:i!==0&&(!E||E==="center"||E==="start"||E==="end")}),d.children?{default:d.children.tab}:void 0)))):k.map((d,i)=>(u.value.push(d.props.name),Ce(i!==0&&!E?Oe(d):d))),!p&&f&&A?je(f,(P?y.length:k.length)!==0):null,E?null:g("div",{class:`${e}-tabs-scroll-padding`,style:{width:`${this.tabsPadding}px`}}));return g("div",{ref:"tabsElRef",class:`${e}-tabs-nav-scroll-content`},A&&f?g(ge,{onResize:this.handleTabsResize},{default:()=>S}):S,A?g("div",{class:`${e}-tabs-pad`}):null,A?null:g("div",{ref:"barElRef",class:`${e}-tabs-bar`}))},O=V?"top":l;return g("div",{class:[`${e}-tabs`,this.themeClass,`${e}-tabs--${n}-type`,`${e}-tabs--${m}-size`,E&&`${e}-tabs--flex`,`${e}-tabs--${O}`],style:this.cssVars},g("div",{class:[`${e}-tabs-nav--${n}-type`,`${e}-tabs-nav--${O}`,`${e}-tabs-nav`]},Be(R,S=>S&&g("div",{class:`${e}-tabs-nav__prefix`},S)),V?g(ge,{onResize:this.handleSegmentResize},{default:()=>g("div",{class:`${e}-tabs-rail`,ref:"tabsElRef"},g("div",{class:`${e}-tabs-capsule`,ref:"segmentCapsuleElRef"},g("div",{class:`${e}-tabs-wrapper`},g("div",{class:`${e}-tabs-tab`}))),P?y.map((S,d)=>(u.value.push(S.props.name),g(we,Object.assign({},S.props,{internalCreatedByPane:!0,internalLeftPadded:d!==0}),S.children?{default:S.children.tab}:void 0))):k.map((S,d)=>(u.value.push(S.props.name),d===0?S:Oe(S))))}):g(ge,{onResize:this.handleNavResize},{default:()=>g("div",{class:`${e}-tabs-nav-scroll-wrapper`,ref:"scrollWrapperElRef"},["top","bottom"].includes(O)?g(ca,{ref:"xScrollInstRef",onScroll:this.handleScroll},{default:Y}):g("div",{class:`${e}-tabs-nav-y-scroll`,onScroll:this.handleScroll,ref:"yScrollElRef"},Y()))}),p&&f&&A?je(f,!0):null,Be(T,S=>S&&g("div",{class:`${e}-tabs-nav__suffix`},S))),P&&(this.animated&&(O==="top"||O==="bottom")?g("div",{ref:"tabsPaneWrapperRef",style:x,class:[`${e}-tabs-pane-wrapper`,_]},Ee(y,this.mergedValue,this.renderedNames,this.onAnimationBeforeLeave,this.onAnimationEnter,this.onAnimationAfterEnter,this.animationDirection)):Ee(y,this.mergedValue,this.renderedNames)))}});function Ee(e,n,l,p,f,m,u){const C=[];return e.forEach(_=>{const{name:x,displayDirective:h,"display-directive":R}=_.props,T=k=>h===k||R===k,y=n===x;if(_.key!==void 0&&(_.key=x),y||T("show")||T("show:lazy")&&l.has(x)){l.has(x)||l.add(x);const k=!T("if");C.push(k?Kt(_,[[Yt,y]]):_)}}),u?g(qt,{name:`${u}-transition`,onBeforeLeave:p,onEnter:f,onAfterEnter:m},{default:()=>C}):C}function je(e,n){return g(we,{ref:"addTabInstRef",key:"__addable",name:"__addable",internalCreatedByPane:!0,internalAddable:!0,internalLeftPadded:n,disabled:typeof e=="object"&&e.disabled})}function Oe(e){const n=Jt(e);return n.props?n.props.internalLeftPadded=!0:n.props={internalLeftPadded:!0},n}function Ce(e){return Array.isArray(e.dynamicProps)?e.dynamicProps.includes("internalLeftPadded")||e.dynamicProps.push("internalLeftPadded"):e.dynamicProps=["internalLeftPadded"],e}const xa=12e3,Ie="185811";function ya(e){const n=_e(e);return e.kind==="gdstudio"?`${n}/api.php`:`${n}/?server=netease&type=url&id=${Ie}`}function Ve(e){return _e(e)||Qt(e.kind)}function Sa(e){return typeof e=="string"&&/^https?:\/\//i.test(e)}function Ca(e,n){if(e<200||e>=400)return{ok:!1,error:`HTTP ${e}`};const l=n.trim();if(!l)return{ok:!1,error:"空响应"};if(/^https?:\/\//i.test(l))return{ok:!0};try{const p=JSON.parse(l);return Sa(p==null?void 0:p.url)?{ok:!0}:{ok:!1,error:"无播放地址"}}catch{return/<html/i.test(l)?{ok:!1,error:"返回了网页"}:{ok:!1,error:"无法解析"}}}async function wa(e){const n=await Zt(e,Number(Ie),320);if(n)return{ok:!0,body:n.slice(0,180)};const l=await ea(e,{name:"海阔天空",artist:"Beyond"},320);return l?{ok:!0,body:l.slice(0,180)}:{ok:!1,error:"无播放地址"}}async function _a(e){const n=_e(e),l=performance.now(),p=new AbortController,f=setTimeout(()=>p.abort(),xa);try{if(e.kind==="gdstudio"){const R=await wa(n),T=Math.round(performance.now()-l);return R.ok?{ok:!0,latencyMs:T,body:R.body}:{ok:!1,latencyMs:T,error:R.error,body:R.body}}const m=ya(e),u=await fetch(m,{method:"GET",cache:"no-store",redirect:"manual",signal:p.signal}),C=Math.round(performance.now()-l),_=u.headers.get("Location");if(_&&/^https?:\/\//i.test(_))return{ok:!0,latencyMs:C,body:_.slice(0,500)};const x=(await u.text()).slice(0,500),h=Ca(u.status,x);return h.ok?{ok:!0,latencyMs:C,body:x||void 0}:{ok:!1,latencyMs:C,error:h.error,body:x||void 0}}catch(m){let u=m instanceof Error?m.message:String(m);return/failed to fetch|networkerror|load failed/i.test(u)&&(u="Failed to fetch（多半未放行或源已挂）"),{ok:!1,latencyMs:Math.round(performance.now()-l),error:u}}finally{clearTimeout(f)}}const Pa={class:"monitor-page"},Ta={class:"monitor-page__header"},Ra=["disabled"],ka={class:"monitor-cards"},za={class:"monitor-card__top"},$a={class:"monitor-card__title"},La={class:"monitor-card__url"},Ba={class:"monitor-card__badge"},Aa={class:"monitor-card__meta"},Ma={key:0,class:"monitor-card__error"},Wa={key:1,class:"monitor-card__body"},Ea=["disabled","onClick"],ja={class:"monitor-cards"},Oa={class:"monitor-card__top"},Va={class:"monitor-card__title"},Ga={class:"monitor-card__url"},Ha={class:"monitor-card__badge"},Da={class:"monitor-field"},Fa={class:"monitor-card__meta"},Ia={key:0,class:"monitor-card__error"},Na={key:1,class:"monitor-card__body"},Ua={class:"monitor-card__actions"},Xa=["disabled","onClick"],Ka=["onClick"],Ya=["disabled","onClick"],qa=["disabled","onClick"],Ja=["onClick"],Qa=["disabled"],I="play-sources",Za=re({__name:"MonitorPage",setup(e){const n=[{id:"onij.fun",label:"onij.fun",cards:[{id:"ping",title:"Ping",url:"https://onij.fun/go/ping"}]},{id:"172.26.250.1",label:"172.26.250.1",cards:[{id:"uvr",title:"UVR",url:"http://172.26.250.1:5555/health"},{id:"whisper",title:"Whisper",url:"http://172.26.250.1:5555/api/v1/whisper/health"}]}],l=[{label:"Meting",value:"meting"},{label:"GD Studio",value:"gdstudio"}],p=ta(),f=L(n[0].id),m=L(!1),u=aa({}),C=L(""),_=L("meting"),x=L("");function h(d,i){return`${d}::${i}`}function R(d,i){var r;return((r=u[h(d,i)])==null?void 0:r.status)??"idle"}function T(d,i){var r;return((r=u[h(d,i)])==null?void 0:r.error)??""}function y(d,i){var r;return((r=u[h(d,i)])==null?void 0:r.body)??""}function k(d,i){var w;const r=(w=u[h(d,i)])==null?void 0:w.latencyMs;return typeof r=="number"?`${r} ms`:"—"}function P(d,i){var w;const r=(w=u[h(d,i)])==null?void 0:w.checkedAt;return r?new Date(r).toLocaleTimeString():"—"}function A(d){switch(d){case"loading":return"检测中";case"ok":return"正常";case"fail":return"异常";default:return"未检测"}}async function V(d,i){const r=h(d,i.id);u[r]={status:"loading"};const w=performance.now();try{const B=await fetch(i.url,{method:"GET",cache:"no-store",signal:AbortSignal.timeout(8e3)}),G=Math.round(performance.now()-w),te=(await B.text()).slice(0,500);if(!B.ok){u[r]={status:"fail",latencyMs:G,checkedAt:Date.now(),error:`HTTP ${B.status}`,body:te||void 0};return}u[r]={status:"ok",latencyMs:G,checkedAt:Date.now(),body:te||void 0}}catch(B){let G=B instanceof Error?B.message:String(B);/failed to fetch|networkerror|load failed/i.test(G)&&(G="Failed to fetch（多半未放行）"),u[r]={status:"fail",latencyMs:Math.round(performance.now()-w),checkedAt:Date.now(),error:G}}}async function E(d){const i=h(I,d.id);u[i]={status:"loading"};const r=await _a(d);u[i]={status:r.ok?"ok":"fail",latencyMs:r.latencyMs,checkedAt:Date.now(),error:r.error,body:r.body}}function Y(d){return d==="gdstudio"?"GD Studio":"Meting"}function O(){const d=x.value.trim();d&&(p.addSource({title:C.value,kind:_.value,origin:d}),C.value="",x.value="")}async function S(){m.value=!0;try{if(f.value===I){await Promise.all(p.sources.map(i=>E(i)));return}const d=n.find(i=>i.id===f.value);if(!d)return;await Promise.all(d.cards.map(i=>V(d.id,i)))}finally{m.value=!1}}return be(f,()=>{S()}),De(()=>{S()}),(d,i)=>(F(),ra(ia,null,{default:le(()=>[b("div",Pa,[b("div",Ta,[i[4]||(i[4]=b("h1",{class:"browse-section-title border-l-4 border-blue-500 pl-2"},"监控",-1)),b("button",{type:"button",class:"monitor-btn",disabled:m.value,onClick:S},W(m.value?"检测中…":"刷新当前标签"),9,Ra)]),q(j(ma),{value:f.value,"onUpdate:value":i[3]||(i[3]=r=>f.value=r),type:"line",animated:"",class:"monitor-tabs"},{default:le(()=>[(F(),X(de,null,me(n,r=>q(j(We),{key:r.id,name:r.id,tab:r.label},{default:le(()=>[b("div",ka,[(F(!0),X(de,null,me(r.cards,w=>(F(),X("article",{key:w.id,class:Me(["monitor-card",`monitor-card--${R(r.id,w.id)}`])},[b("div",za,[b("div",null,[b("h2",$a,W(w.title),1),b("p",La,W(w.url),1)]),b("span",Ba,W(A(R(r.id,w.id))),1)]),b("dl",Aa,[b("div",null,[i[5]||(i[5]=b("dt",null,"延迟",-1)),b("dd",null,W(k(r.id,w.id)),1)]),b("div",null,[i[6]||(i[6]=b("dt",null,"最近检测",-1)),b("dd",null,W(P(r.id,w.id)),1)])]),T(r.id,w.id)?(F(),X("p",Ma,W(T(r.id,w.id)),1)):y(r.id,w.id)?(F(),X("pre",Wa,W(y(r.id,w.id)),1)):xe("",!0),b("button",{type:"button",class:"monitor-btn monitor-btn--ghost",disabled:m.value,onClick:B=>V(r.id,w)}," 单独检测 ",8,Ea)],2))),128))])]),_:2},1032,["name","tab"])),64)),q(j(We),{name:I,tab:"音乐源"},{default:le(()=>[i[11]||(i[11]=b("p",{class:"monitor-hint"},"取链按卡片从上到下尝试。可改地址、开关、排序；内建源只能停用不能删。",-1)),b("div",ja,[(F(!0),X(de,null,me(j(p).sources,(r,w)=>(F(),X("article",{key:r.id,class:Me(["monitor-card",[`monitor-card--${R(I,r.id)}`,{"monitor-card--off":!r.enabled}]])},[b("div",Oa,[b("div",null,[b("h2",Va,W(r.title),1),b("p",Ga,W(Y(r.kind))+" · "+W(j(Ve)(r)),1)]),b("span",Ha,W(A(R(I,r.id))),1)]),b("label",Da,[i[7]||(i[7]=b("span",null,"源地址",-1)),q(j(ye),{value:r.origin,placeholder:j(Ve)(r),size:"small","onUpdate:value":B=>j(p).patchSource(r.id,{origin:B})},null,8,["value","placeholder","onUpdate:value"])]),b("dl",Fa,[b("div",null,[i[8]||(i[8]=b("dt",null,"延迟",-1)),b("dd",null,W(k(I,r.id)),1)]),b("div",null,[i[9]||(i[9]=b("dt",null,"最近检测",-1)),b("dd",null,W(P(I,r.id)),1)])]),T(I,r.id)?(F(),X("p",Ia,W(T(I,r.id)),1)):y(I,r.id)?(F(),X("pre",Na,W(y(I,r.id)),1)):xe("",!0),b("div",Ua,[b("button",{type:"button",class:"monitor-btn monitor-btn--ghost",disabled:m.value,onClick:B=>E(r)}," 单独检测 ",8,Xa),b("button",{type:"button",class:"monitor-btn monitor-btn--ghost",onClick:B=>j(p).patchSource(r.id,{enabled:!r.enabled})},W(r.enabled?"停用":"启用"),9,Ka),b("button",{type:"button",class:"monitor-btn monitor-btn--ghost",disabled:w===0,onClick:B=>j(p).moveSource(r.id,-1)}," 上移 ",8,Ya),b("button",{type:"button",class:"monitor-btn monitor-btn--ghost",disabled:w===j(p).sources.length-1,onClick:B=>j(p).moveSource(r.id,1)}," 下移 ",8,qa),r.builtin?xe("",!0):(F(),X("button",{key:0,type:"button",class:"monitor-btn monitor-btn--ghost",onClick:B=>j(p).removeSource(r.id)}," 删除 ",8,Ja))])],2))),128))]),b("form",{class:"monitor-add",onSubmit:oa(O,["prevent"])},[i[10]||(i[10]=b("h3",{class:"monitor-add__title"},"添加源",-1)),q(j(ye),{value:C.value,"onUpdate:value":i[0]||(i[0]=r=>C.value=r),size:"small",placeholder:"名称"},null,8,["value"]),q(j(la),{value:_.value,"onUpdate:value":i[1]||(i[1]=r=>_.value=r),size:"small",options:l},null,8,["value"]),q(j(ye),{value:x.value,"onUpdate:value":i[2]||(i[2]=r=>x.value=r),size:"small",placeholder:"https://example.com/meting"},null,8,["value"]),b("button",{type:"submit",class:"monitor-btn",disabled:!x.value.trim()},"添加",8,Qa)],32)]),_:1})]),_:1},8,["value"])])]),_:1}))}}),rr=na(Za,[["__scopeId","data-v-df7a86c2"]]);export{rr as default};
