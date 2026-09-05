import{P as yt}from"./PageContent-DlI9g3oB.js";import{d as Y,aa as f,r as L,ad as Se,ae as Ct,af as St,bY as wt,bZ as Tt,aA as Pt,aX as _t,a5 as Be,bp as Rt,ac as zt,aH as $t,F as pe,az as Lt,aU as Bt,c as te,b8 as Wt,as as r,aK as d,at as P,au as z,aL as At,bH as ie,aN as we,ab as se,aB as kt,aC as We,bh as Te,b2 as Et,w as ae,o as Ae,a9 as jt,ai as M,b_ as Mt,a$ as Vt,aD as V,aR as Q,aE as Ht,aj as Ot,A as Gt,bf as Dt,b$ as Ft,bx as It,am as le,aW as ee,a as Nt,e as Ut,f as de,g as R,n as Pe,t as G,h as K,j as _e,y as Re,l as U,I as Xt,i as Kt,_ as Yt}from"./index-BzjSw0KY.js";import{A as qt}from"./Add-B40CM-l4.js";const Jt=Se(".v-x-scroll",{overflow:"auto",scrollbarWidth:"none"},[Se("&::-webkit-scrollbar",{width:0,height:0})]),Zt=Y({name:"XScroll",props:{disabled:Boolean,onScroll:Function},setup(){const e=L(null);function i(l){!(l.currentTarget.offsetWidth<l.currentTarget.scrollWidth)||l.deltaY===0||(l.currentTarget.scrollLeft+=l.deltaY+l.deltaX,l.preventDefault())}const c=Ct();return Jt.mount({id:"vueuc/x-scroll",head:!0,anchorMetaName:St,ssr:c}),Object.assign({selfRef:e,handleWheel:i},{scrollTo(...l){var h;(h=e.value)===null||h===void 0||h.scrollTo(...l)}})},render(){return f("div",{ref:"selfRef",onScroll:this.onScroll,onWheel:this.disabled?void 0:this.handleWheel,class:"v-x-scroll"},this.$slots)}});var Qt="Expected a function";function ce(e,i,c){var v=!0,l=!0;if(typeof e!="function")throw new TypeError(Qt);return Tt(c)&&(v="leading"in c?!!c.leading:v,l="trailing"in c?!!c.trailing:l),wt(e,i,{leading:v,maxWait:i,trailing:l})}const ea={tabFontSizeSmall:"14px",tabFontSizeMedium:"14px",tabFontSizeLarge:"16px",tabGapSmallLine:"36px",tabGapMediumLine:"36px",tabGapLargeLine:"36px",tabGapSmallLineVertical:"8px",tabGapMediumLineVertical:"8px",tabGapLargeLineVertical:"8px",tabPaddingSmallLine:"6px 0",tabPaddingMediumLine:"10px 0",tabPaddingLargeLine:"14px 0",tabPaddingVerticalSmallLine:"6px 12px",tabPaddingVerticalMediumLine:"8px 16px",tabPaddingVerticalLargeLine:"10px 20px",tabGapSmallBar:"36px",tabGapMediumBar:"36px",tabGapLargeBar:"36px",tabGapSmallBarVertical:"8px",tabGapMediumBarVertical:"8px",tabGapLargeBarVertical:"8px",tabPaddingSmallBar:"4px 0",tabPaddingMediumBar:"6px 0",tabPaddingLargeBar:"10px 0",tabPaddingVerticalSmallBar:"6px 12px",tabPaddingVerticalMediumBar:"8px 16px",tabPaddingVerticalLargeBar:"10px 20px",tabGapSmallCard:"4px",tabGapMediumCard:"4px",tabGapLargeCard:"4px",tabGapSmallCardVertical:"4px",tabGapMediumCardVertical:"4px",tabGapLargeCardVertical:"4px",tabPaddingSmallCard:"8px 16px",tabPaddingMediumCard:"10px 20px",tabPaddingLargeCard:"12px 24px",tabPaddingSmallSegment:"4px 0",tabPaddingMediumSegment:"6px 0",tabPaddingLargeSegment:"8px 0",tabPaddingVerticalLargeSegment:"0 8px",tabPaddingVerticalSmallCard:"8px 12px",tabPaddingVerticalMediumCard:"10px 16px",tabPaddingVerticalLargeCard:"12px 20px",tabPaddingVerticalSmallSegment:"0 4px",tabPaddingVerticalMediumSegment:"0 6px",tabGapSmallSegment:"0",tabGapMediumSegment:"0",tabGapLargeSegment:"0",tabGapSmallSegmentVertical:"0",tabGapMediumSegmentVertical:"0",tabGapLargeSegmentVertical:"0",panePaddingSmall:"8px 0 0 0",panePaddingMedium:"12px 0 0 0",panePaddingLarge:"16px 0 0 0",closeSize:"18px",closeIconSize:"14px"};function ta(e){const{textColor2:i,primaryColor:c,textColorDisabled:v,closeIconColor:l,closeIconColorHover:h,closeIconColorPressed:y,closeColorHover:w,closeColorPressed:_,tabColor:C,baseColor:x,dividerColor:A,fontWeight:$,textColor1:g,borderRadius:b,fontSize:s,fontWeightStrong:n}=e;return Object.assign(Object.assign({},ea),{colorSegment:C,tabFontSizeCard:s,tabTextColorLine:g,tabTextColorActiveLine:c,tabTextColorHoverLine:c,tabTextColorDisabledLine:v,tabTextColorSegment:g,tabTextColorActiveSegment:i,tabTextColorHoverSegment:i,tabTextColorDisabledSegment:v,tabTextColorBar:g,tabTextColorActiveBar:c,tabTextColorHoverBar:c,tabTextColorDisabledBar:v,tabTextColorCard:g,tabTextColorHoverCard:g,tabTextColorActiveCard:c,tabTextColorDisabledCard:v,barColor:c,closeIconColor:l,closeIconColorHover:h,closeIconColorPressed:y,closeColorHover:w,closeColorPressed:_,closeBorderRadius:b,tabColor:C,tabColorSegment:x,tabBorderColor:A,tabFontWeightActive:$,tabFontWeight:$,tabBorderRadius:b,paneTextColor:i,fontWeightStrong:n})}const aa={common:Pt,self:ta},ue=_t("n-tabs"),ke={tab:[String,Number,Object,Function],name:{type:[String,Number],required:!0},disabled:Boolean,displayDirective:{type:String,default:"if"},closable:{type:Boolean,default:void 0},tabProps:Object,label:[String,Number,Object,Function]},ra=Y({__TAB_PANE__:!0,name:"TabPane",alias:["TabPanel"],props:ke,slots:Object,setup(e){const i=Be(ue,null);return i||Rt("tab-pane","`n-tab-pane` must be placed inside `n-tabs`."),{style:i.paneStyleRef,class:i.paneClassRef,mergedClsPrefix:i.mergedClsPrefixRef}},render(){return f("div",{class:[`${this.mergedClsPrefix}-tab-pane`,this.class],style:this.style},this.$slots)}}),oa=Object.assign({internalLeftPadded:Boolean,internalAddable:Boolean,internalCreatedByPane:Boolean},Wt(ke,["displayDirective"])),fe=Y({__TAB__:!0,inheritAttrs:!1,name:"Tab",props:oa,setup(e){const{mergedClsPrefixRef:i,valueRef:c,typeRef:v,closableRef:l,tabStyleRef:h,addTabStyleRef:y,tabClassRef:w,addTabClassRef:_,tabChangeIdRef:C,onBeforeLeaveRef:x,triggerRef:A,handleAdd:$,activateTab:g,handleClose:b}=Be(ue);return{trigger:A,mergedClosable:te(()=>{if(e.internalAddable)return!1;const{closable:s}=e;return s===void 0?l.value:s}),style:h,addStyle:y,tabClass:w,addTabClass:_,clsPrefix:i,value:c,type:v,handleClose(s){s.stopPropagation(),!e.disabled&&b(e.name)},activateTab(){if(e.disabled)return;if(e.internalAddable){$();return}const{name:s}=e,n=++C.id;if(s!==c.value){const{value:u}=x;u?Promise.resolve(u(e.name,c.value)).then(S=>{S&&C.id===n&&g(s)}):g(s)}}}},render(){const{internalAddable:e,clsPrefix:i,name:c,disabled:v,label:l,tab:h,value:y,mergedClosable:w,trigger:_,$slots:{default:C}}=this,x=l??h;return f("div",{class:`${i}-tabs-tab-wrapper`},this.internalLeftPadded?f("div",{class:`${i}-tabs-tab-pad`}):null,f("div",Object.assign({key:c,"data-name":c,"data-disabled":v?!0:void 0},zt({class:[`${i}-tabs-tab`,y===c&&`${i}-tabs-tab--active`,v&&`${i}-tabs-tab--disabled`,w&&`${i}-tabs-tab--closable`,e&&`${i}-tabs-tab--addable`,e?this.addTabClass:this.tabClass],onClick:_==="click"?this.activateTab:void 0,onMouseenter:_==="hover"?this.activateTab:void 0,style:e?this.addStyle:this.style},this.internalCreatedByPane?this.tabProps||{}:this.$attrs)),f("span",{class:`${i}-tabs-tab__label`},e?f(pe,null,f("div",{class:`${i}-tabs-tab__height-placeholder`}," "),f(Lt,{clsPrefix:i},{default:()=>f(qt,null)})):C?C():typeof x=="object"?x:$t(x??c)),w&&this.type==="card"?f(Bt,{clsPrefix:i,class:`${i}-tabs-tab__close`,onClick:this.handleClose,disabled:v}):null))}}),na=r("tabs",`
 box-sizing: border-box;
 width: 100%;
 display: flex;
 flex-direction: column;
 transition:
 background-color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
`,[d("segment-type",[r("tabs-rail",[P("&.transition-disabled",[r("tabs-capsule",`
 transition: none;
 `)])])]),d("top",[r("tab-pane",`
 padding: var(--n-pane-padding-top) var(--n-pane-padding-right) var(--n-pane-padding-bottom) var(--n-pane-padding-left);
 `)]),d("left",[r("tab-pane",`
 padding: var(--n-pane-padding-right) var(--n-pane-padding-bottom) var(--n-pane-padding-left) var(--n-pane-padding-top);
 `)]),d("left, right",`
 flex-direction: row;
 `,[r("tabs-bar",`
 width: 2px;
 right: 0;
 transition:
 top .2s var(--n-bezier),
 max-height .2s var(--n-bezier),
 background-color .3s var(--n-bezier);
 `),r("tabs-tab",`
 padding: var(--n-tab-padding-vertical); 
 `)]),d("right",`
 flex-direction: row-reverse;
 `,[r("tab-pane",`
 padding: var(--n-pane-padding-left) var(--n-pane-padding-top) var(--n-pane-padding-right) var(--n-pane-padding-bottom);
 `),r("tabs-bar",`
 left: 0;
 `)]),d("bottom",`
 flex-direction: column-reverse;
 justify-content: flex-end;
 `,[r("tab-pane",`
 padding: var(--n-pane-padding-bottom) var(--n-pane-padding-right) var(--n-pane-padding-top) var(--n-pane-padding-left);
 `),r("tabs-bar",`
 top: 0;
 `)]),r("tabs-rail",`
 position: relative;
 padding: 3px;
 border-radius: var(--n-tab-border-radius);
 width: 100%;
 background-color: var(--n-color-segment);
 transition: background-color .3s var(--n-bezier);
 display: flex;
 align-items: center;
 `,[r("tabs-capsule",`
 border-radius: var(--n-tab-border-radius);
 position: absolute;
 pointer-events: none;
 background-color: var(--n-tab-color-segment);
 box-shadow: 0 1px 3px 0 rgba(0, 0, 0, .08);
 transition: transform 0.3s var(--n-bezier);
 `),r("tabs-tab-wrapper",`
 flex-basis: 0;
 flex-grow: 1;
 display: flex;
 align-items: center;
 justify-content: center;
 `,[r("tabs-tab",`
 overflow: hidden;
 border-radius: var(--n-tab-border-radius);
 width: 100%;
 display: flex;
 align-items: center;
 justify-content: center;
 `,[d("active",`
 font-weight: var(--n-font-weight-strong);
 color: var(--n-tab-text-color-active);
 `),P("&:hover",`
 color: var(--n-tab-text-color-hover);
 `)])])]),d("flex",[r("tabs-nav",`
 width: 100%;
 position: relative;
 `,[r("tabs-wrapper",`
 width: 100%;
 `,[r("tabs-tab",`
 margin-right: 0;
 `)])])]),r("tabs-nav",`
 box-sizing: border-box;
 line-height: 1.5;
 display: flex;
 transition: border-color .3s var(--n-bezier);
 `,[z("prefix, suffix",`
 display: flex;
 align-items: center;
 `),z("prefix","padding-right: 16px;"),z("suffix","padding-left: 16px;")]),d("top, bottom",[r("tabs-nav-scroll-wrapper",[P("&::before",`
 top: 0;
 bottom: 0;
 left: 0;
 width: 20px;
 `),P("&::after",`
 top: 0;
 bottom: 0;
 right: 0;
 width: 20px;
 `),d("shadow-start",[P("&::before",`
 box-shadow: inset 10px 0 8px -8px rgba(0, 0, 0, .12);
 `)]),d("shadow-end",[P("&::after",`
 box-shadow: inset -10px 0 8px -8px rgba(0, 0, 0, .12);
 `)])])]),d("left, right",[r("tabs-nav-scroll-content",`
 flex-direction: column;
 `),r("tabs-nav-scroll-wrapper",[P("&::before",`
 top: 0;
 left: 0;
 right: 0;
 height: 20px;
 `),P("&::after",`
 bottom: 0;
 left: 0;
 right: 0;
 height: 20px;
 `),d("shadow-start",[P("&::before",`
 box-shadow: inset 0 10px 8px -8px rgba(0, 0, 0, .12);
 `)]),d("shadow-end",[P("&::after",`
 box-shadow: inset 0 -10px 8px -8px rgba(0, 0, 0, .12);
 `)])])]),r("tabs-nav-scroll-wrapper",`
 flex: 1;
 position: relative;
 overflow: hidden;
 `,[r("tabs-nav-y-scroll",`
 height: 100%;
 width: 100%;
 overflow-y: auto; 
 scrollbar-width: none;
 `,[P("&::-webkit-scrollbar, &::-webkit-scrollbar-track-piece, &::-webkit-scrollbar-thumb",`
 width: 0;
 height: 0;
 display: none;
 `)]),P("&::before, &::after",`
 transition: box-shadow .3s var(--n-bezier);
 pointer-events: none;
 content: "";
 position: absolute;
 z-index: 1;
 `)]),r("tabs-nav-scroll-content",`
 display: flex;
 position: relative;
 min-width: 100%;
 min-height: 100%;
 width: fit-content;
 box-sizing: border-box;
 `),r("tabs-wrapper",`
 display: inline-flex;
 flex-wrap: nowrap;
 position: relative;
 `),r("tabs-tab-wrapper",`
 display: flex;
 flex-wrap: nowrap;
 flex-shrink: 0;
 flex-grow: 0;
 `),r("tabs-tab",`
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
 `,[d("disabled",{cursor:"not-allowed"}),z("close",`
 margin-left: 6px;
 transition:
 background-color .3s var(--n-bezier),
 color .3s var(--n-bezier);
 `),z("label",`
 display: flex;
 align-items: center;
 z-index: 1;
 `)]),r("tabs-bar",`
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
 `,[P("&.transition-disabled",`
 transition: none;
 `),d("disabled",`
 background-color: var(--n-tab-text-color-disabled)
 `)]),r("tabs-pane-wrapper",`
 position: relative;
 overflow: hidden;
 transition: max-height .2s var(--n-bezier);
 `),r("tab-pane",`
 color: var(--n-pane-text-color);
 width: 100%;
 transition:
 color .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 opacity .2s var(--n-bezier);
 left: 0;
 right: 0;
 top: 0;
 `,[P("&.next-transition-leave-active, &.prev-transition-leave-active, &.next-transition-enter-active, &.prev-transition-enter-active",`
 transition:
 color .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 transform .2s var(--n-bezier),
 opacity .2s var(--n-bezier);
 `),P("&.next-transition-leave-active, &.prev-transition-leave-active",`
 position: absolute;
 `),P("&.next-transition-enter-from, &.prev-transition-leave-to",`
 transform: translateX(32px);
 opacity: 0;
 `),P("&.next-transition-leave-to, &.prev-transition-enter-from",`
 transform: translateX(-32px);
 opacity: 0;
 `),P("&.next-transition-leave-from, &.next-transition-enter-to, &.prev-transition-leave-from, &.prev-transition-enter-to",`
 transform: translateX(0);
 opacity: 1;
 `)]),r("tabs-tab-pad",`
 box-sizing: border-box;
 width: var(--n-tab-gap);
 flex-grow: 0;
 flex-shrink: 0;
 `),d("line-type, bar-type",[r("tabs-tab",`
 font-weight: var(--n-tab-font-weight);
 box-sizing: border-box;
 vertical-align: bottom;
 `,[P("&:hover",{color:"var(--n-tab-text-color-hover)"}),d("active",`
 color: var(--n-tab-text-color-active);
 font-weight: var(--n-tab-font-weight-active);
 `),d("disabled",{color:"var(--n-tab-text-color-disabled)"})])]),r("tabs-nav",[d("line-type",[d("top",[z("prefix, suffix",`
 border-bottom: 1px solid var(--n-tab-border-color);
 `),r("tabs-nav-scroll-content",`
 border-bottom: 1px solid var(--n-tab-border-color);
 `),r("tabs-bar",`
 bottom: -1px;
 `)]),d("left",[z("prefix, suffix",`
 border-right: 1px solid var(--n-tab-border-color);
 `),r("tabs-nav-scroll-content",`
 border-right: 1px solid var(--n-tab-border-color);
 `),r("tabs-bar",`
 right: -1px;
 `)]),d("right",[z("prefix, suffix",`
 border-left: 1px solid var(--n-tab-border-color);
 `),r("tabs-nav-scroll-content",`
 border-left: 1px solid var(--n-tab-border-color);
 `),r("tabs-bar",`
 left: -1px;
 `)]),d("bottom",[z("prefix, suffix",`
 border-top: 1px solid var(--n-tab-border-color);
 `),r("tabs-nav-scroll-content",`
 border-top: 1px solid var(--n-tab-border-color);
 `),r("tabs-bar",`
 top: -1px;
 `)]),z("prefix, suffix",`
 transition: border-color .3s var(--n-bezier);
 `),r("tabs-nav-scroll-content",`
 transition: border-color .3s var(--n-bezier);
 `),r("tabs-bar",`
 border-radius: 0;
 `)]),d("card-type",[z("prefix, suffix",`
 transition: border-color .3s var(--n-bezier);
 `),r("tabs-pad",`
 flex-grow: 1;
 transition: border-color .3s var(--n-bezier);
 `),r("tabs-tab-pad",`
 transition: border-color .3s var(--n-bezier);
 `),r("tabs-tab",`
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
 `,[d("addable",`
 padding-left: 8px;
 padding-right: 8px;
 font-size: 16px;
 justify-content: center;
 `,[z("height-placeholder",`
 width: 0;
 font-size: var(--n-tab-font-size);
 `),At("disabled",[P("&:hover",`
 color: var(--n-tab-text-color-hover);
 `)])]),d("closable","padding-right: 8px;"),d("active",`
 background-color: #0000;
 font-weight: var(--n-tab-font-weight-active);
 color: var(--n-tab-text-color-active);
 `),d("disabled","color: var(--n-tab-text-color-disabled);")])]),d("left, right",`
 flex-direction: column; 
 `,[z("prefix, suffix",`
 padding: var(--n-tab-padding-vertical);
 `),r("tabs-wrapper",`
 flex-direction: column;
 `),r("tabs-tab-wrapper",`
 flex-direction: column;
 `,[r("tabs-tab-pad",`
 height: var(--n-tab-gap-vertical);
 width: 100%;
 `)])]),d("top",[d("card-type",[r("tabs-scroll-padding","border-bottom: 1px solid var(--n-tab-border-color);"),z("prefix, suffix",`
 border-bottom: 1px solid var(--n-tab-border-color);
 `),r("tabs-tab",`
 border-top-left-radius: var(--n-tab-border-radius);
 border-top-right-radius: var(--n-tab-border-radius);
 `,[d("active",`
 border-bottom: 1px solid #0000;
 `)]),r("tabs-tab-pad",`
 border-bottom: 1px solid var(--n-tab-border-color);
 `),r("tabs-pad",`
 border-bottom: 1px solid var(--n-tab-border-color);
 `)])]),d("left",[d("card-type",[r("tabs-scroll-padding","border-right: 1px solid var(--n-tab-border-color);"),z("prefix, suffix",`
 border-right: 1px solid var(--n-tab-border-color);
 `),r("tabs-tab",`
 border-top-left-radius: var(--n-tab-border-radius);
 border-bottom-left-radius: var(--n-tab-border-radius);
 `,[d("active",`
 border-right: 1px solid #0000;
 `)]),r("tabs-tab-pad",`
 border-right: 1px solid var(--n-tab-border-color);
 `),r("tabs-pad",`
 border-right: 1px solid var(--n-tab-border-color);
 `)])]),d("right",[d("card-type",[r("tabs-scroll-padding","border-left: 1px solid var(--n-tab-border-color);"),z("prefix, suffix",`
 border-left: 1px solid var(--n-tab-border-color);
 `),r("tabs-tab",`
 border-top-right-radius: var(--n-tab-border-radius);
 border-bottom-right-radius: var(--n-tab-border-radius);
 `,[d("active",`
 border-left: 1px solid #0000;
 `)]),r("tabs-tab-pad",`
 border-left: 1px solid var(--n-tab-border-color);
 `),r("tabs-pad",`
 border-left: 1px solid var(--n-tab-border-color);
 `)])]),d("bottom",[d("card-type",[r("tabs-scroll-padding","border-top: 1px solid var(--n-tab-border-color);"),z("prefix, suffix",`
 border-top: 1px solid var(--n-tab-border-color);
 `),r("tabs-tab",`
 border-bottom-left-radius: var(--n-tab-border-radius);
 border-bottom-right-radius: var(--n-tab-border-radius);
 `,[d("active",`
 border-top: 1px solid #0000;
 `)]),r("tabs-tab-pad",`
 border-top: 1px solid var(--n-tab-border-color);
 `),r("tabs-pad",`
 border-top: 1px solid var(--n-tab-border-color);
 `)])])])]),ia=Object.assign(Object.assign({},We.props),{value:[String,Number],defaultValue:[String,Number],trigger:{type:String,default:"click"},type:{type:String,default:"bar"},closable:Boolean,justifyContent:String,size:{type:String,default:"medium"},placement:{type:String,default:"top"},tabStyle:[String,Object],tabClass:String,addTabStyle:[String,Object],addTabClass:String,barWidth:Number,paneClass:String,paneStyle:[String,Object],paneWrapperClass:String,paneWrapperStyle:[String,Object],addable:[Boolean,Object],tabsPadding:{type:Number,default:0},animated:Boolean,onBeforeLeave:Function,onAdd:Function,"onUpdate:value":[Function,Array],onUpdateValue:[Function,Array],onClose:[Function,Array],labelSize:String,activeName:[String,Number],onActiveNameChange:[Function,Array]}),sa=Y({name:"Tabs",props:ia,slots:Object,setup(e,{slots:i}){var c,v,l,h;const{mergedClsPrefixRef:y,inlineThemeDisabled:w}=kt(e),_=We("Tabs","-tabs",na,aa,e,y),C=L(null),x=L(null),A=L(null),$=L(null),g=L(null),b=L(null),s=L(!0),n=L(!0),u=Te(e,["labelSize","size"]),S=Te(e,["activeName","value"]),k=L((v=(c=S.value)!==null&&c!==void 0?c:e.defaultValue)!==null&&v!==void 0?v:i.default?(h=(l=ie(i.default())[0])===null||l===void 0?void 0:l.props)===null||h===void 0?void 0:h.name:null),B=Et(S,k),m={id:0},W=te(()=>{if(!(!e.justifyContent||e.type==="card"))return{display:"flex",justifyContent:e.justifyContent}});ae(B,()=>{m.id=0,q(),ge()});function H(){var t;const{value:a}=B;return a===null?null:(t=C.value)===null||t===void 0?void 0:t.querySelector(`[data-name="${a}"]`)}function Ee(t){if(e.type==="card")return;const{value:a}=x;if(!a)return;const o=a.style.opacity==="0";if(t){const p=`${y.value}-tabs-bar--disabled`,{barWidth:T,placement:E}=e;if(t.dataset.disabled==="true"?a.classList.add(p):a.classList.remove(p),["top","bottom"].includes(E)){if(ve(["top","maxHeight","height"]),typeof T=="number"&&t.offsetWidth>=T){const j=Math.floor((t.offsetWidth-T)/2)+t.offsetLeft;a.style.left=`${j}px`,a.style.maxWidth=`${T}px`}else a.style.left=`${t.offsetLeft}px`,a.style.maxWidth=`${t.offsetWidth}px`;a.style.width="8192px",o&&(a.style.transition="none"),a.offsetWidth,o&&(a.style.transition="",a.style.opacity="1")}else{if(ve(["left","maxWidth","width"]),typeof T=="number"&&t.offsetHeight>=T){const j=Math.floor((t.offsetHeight-T)/2)+t.offsetTop;a.style.top=`${j}px`,a.style.maxHeight=`${T}px`}else a.style.top=`${t.offsetTop}px`,a.style.maxHeight=`${t.offsetHeight}px`;a.style.height="8192px",o&&(a.style.transition="none"),a.offsetHeight,o&&(a.style.transition="",a.style.opacity="1")}}}function je(){if(e.type==="card")return;const{value:t}=x;t&&(t.style.opacity="0")}function ve(t){const{value:a}=x;if(a)for(const o of t)a.style[o]=""}function q(){if(e.type==="card")return;const t=H();t?Ee(t):je()}function ge(){var t;const a=(t=g.value)===null||t===void 0?void 0:t.$el;if(!a)return;const o=H();if(!o)return;const{scrollLeft:p,offsetWidth:T}=a,{offsetLeft:E,offsetWidth:j}=o;p>E?a.scrollTo({top:0,left:E,behavior:"smooth"}):E+j>p+T&&a.scrollTo({top:0,left:E+j-T,behavior:"smooth"})}const J=L(null);let re=0,O=null;function Me(t){const a=J.value;if(a){re=t.getBoundingClientRect().height;const o=`${re}px`,p=()=>{a.style.height=o,a.style.maxHeight=o};O?(p(),O(),O=null):O=p}}function Ve(t){const a=J.value;if(a){const o=t.getBoundingClientRect().height,p=()=>{document.body.offsetHeight,a.style.maxHeight=`${o}px`,a.style.height=`${Math.max(re,o)}px`};O?(O(),O=null,p()):O=p}}function He(){const t=J.value;if(t){t.style.maxHeight="",t.style.height="";const{paneWrapperStyle:a}=e;if(typeof a=="string")t.style.cssText=a;else if(a){const{maxHeight:o,height:p}=a;o!==void 0&&(t.style.maxHeight=o),p!==void 0&&(t.style.height=p)}}}const he={value:[]},xe=L("next");function Oe(t){const a=B.value;let o="next";for(const p of he.value){if(p===a)break;if(p===t){o="prev";break}}xe.value=o,Ge(t)}function Ge(t){const{onActiveNameChange:a,onUpdateValue:o,"onUpdate:value":p}=e;a&&ee(a,t),o&&ee(o,t),p&&ee(p,t),k.value=t}function De(t){const{onClose:a}=e;a&&ee(a,t)}function me(){const{value:t}=x;if(!t)return;const a="transition-disabled";t.classList.add(a),q(),t.classList.remove(a)}const D=L(null);function oe({transitionDisabled:t}){const a=C.value;if(!a)return;t&&a.classList.add("transition-disabled");const o=H();o&&D.value&&(D.value.style.width=`${o.offsetWidth}px`,D.value.style.height=`${o.offsetHeight}px`,D.value.style.transform=`translateX(${o.offsetLeft-Ot(getComputedStyle(a).paddingLeft)}px)`,t&&D.value.offsetWidth),t&&a.classList.remove("transition-disabled")}ae([B],()=>{e.type==="segment"&&le(()=>{oe({transitionDisabled:!1})})}),Ae(()=>{e.type==="segment"&&oe({transitionDisabled:!0})});let ye=0;function Fe(t){var a;if(t.contentRect.width===0&&t.contentRect.height===0||ye===t.contentRect.width)return;ye=t.contentRect.width;const{type:o}=e;if((o==="line"||o==="bar")&&me(),o!=="segment"){const{placement:p}=e;ne((p==="top"||p==="bottom"?(a=g.value)===null||a===void 0?void 0:a.$el:b.value)||null)}}const Ie=ce(Fe,64);ae([()=>e.justifyContent,()=>e.size],()=>{le(()=>{const{type:t}=e;(t==="line"||t==="bar")&&me()})});const F=L(!1);function Ne(t){var a;const{target:o,contentRect:{width:p,height:T}}=t,E=o.parentElement.parentElement.offsetWidth,j=o.parentElement.parentElement.offsetHeight,{placement:N}=e;if(!F.value)N==="top"||N==="bottom"?E<p&&(F.value=!0):j<T&&(F.value=!0);else{const{value:X}=$;if(!X)return;N==="top"||N==="bottom"?E-p>X.$el.offsetWidth&&(F.value=!1):j-T>X.$el.offsetHeight&&(F.value=!1)}ne(((a=g.value)===null||a===void 0?void 0:a.$el)||null)}const Ue=ce(Ne,64);function Xe(){const{onAdd:t}=e;t&&t(),le(()=>{const a=H(),{value:o}=g;!a||!o||o.scrollTo({left:a.offsetLeft,top:0,behavior:"smooth"})})}function ne(t){if(!t)return;const{placement:a}=e;if(a==="top"||a==="bottom"){const{scrollLeft:o,scrollWidth:p,offsetWidth:T}=t;s.value=o<=0,n.value=o+T>=p}else{const{scrollTop:o,scrollHeight:p,offsetHeight:T}=t;s.value=o<=0,n.value=o+T>=p}}const Ke=ce(t=>{ne(t.target)},64);jt(ue,{triggerRef:M(e,"trigger"),tabStyleRef:M(e,"tabStyle"),tabClassRef:M(e,"tabClass"),addTabStyleRef:M(e,"addTabStyle"),addTabClassRef:M(e,"addTabClass"),paneClassRef:M(e,"paneClass"),paneStyleRef:M(e,"paneStyle"),mergedClsPrefixRef:y,typeRef:M(e,"type"),closableRef:M(e,"closable"),valueRef:B,tabChangeIdRef:m,onBeforeLeaveRef:M(e,"onBeforeLeave"),activateTab:Oe,handleClose:De,handleAdd:Xe}),Mt(()=>{q(),ge()}),Vt(()=>{const{value:t}=A;if(!t)return;const{value:a}=y,o=`${a}-tabs-nav-scroll-wrapper--shadow-start`,p=`${a}-tabs-nav-scroll-wrapper--shadow-end`;s.value?t.classList.remove(o):t.classList.add(o),n.value?t.classList.remove(p):t.classList.add(p)});const Ye={syncBarPosition:()=>{q()}},qe=()=>{oe({transitionDisabled:!0})},Ce=te(()=>{const{value:t}=u,{type:a}=e,o={card:"Card",bar:"Bar",line:"Line",segment:"Segment"}[a],p=`${t}${o}`,{self:{barColor:T,closeIconColor:E,closeIconColorHover:j,closeIconColorPressed:N,tabColor:X,tabBorderColor:Je,paneTextColor:Ze,tabFontWeight:Qe,tabBorderRadius:et,tabFontWeightActive:tt,colorSegment:at,fontWeightStrong:rt,tabColorSegment:ot,closeSize:nt,closeIconSize:it,closeColorHover:st,closeColorPressed:lt,closeBorderRadius:dt,[V("panePadding",t)]:Z,[V("tabPadding",p)]:ct,[V("tabPaddingVertical",p)]:bt,[V("tabGap",p)]:pt,[V("tabGap",`${p}Vertical`)]:ft,[V("tabTextColor",a)]:ut,[V("tabTextColorActive",a)]:vt,[V("tabTextColorHover",a)]:gt,[V("tabTextColorDisabled",a)]:ht,[V("tabFontSize",t)]:xt},common:{cubicBezierEaseInOut:mt}}=_.value;return{"--n-bezier":mt,"--n-color-segment":at,"--n-bar-color":T,"--n-tab-font-size":xt,"--n-tab-text-color":ut,"--n-tab-text-color-active":vt,"--n-tab-text-color-disabled":ht,"--n-tab-text-color-hover":gt,"--n-pane-text-color":Ze,"--n-tab-border-color":Je,"--n-tab-border-radius":et,"--n-close-size":nt,"--n-close-icon-size":it,"--n-close-color-hover":st,"--n-close-color-pressed":lt,"--n-close-border-radius":dt,"--n-close-icon-color":E,"--n-close-icon-color-hover":j,"--n-close-icon-color-pressed":N,"--n-tab-color":X,"--n-tab-font-weight":Qe,"--n-tab-font-weight-active":tt,"--n-tab-padding":ct,"--n-tab-padding-vertical":bt,"--n-tab-gap":pt,"--n-tab-gap-vertical":ft,"--n-pane-padding-left":Q(Z,"left"),"--n-pane-padding-right":Q(Z,"right"),"--n-pane-padding-top":Q(Z,"top"),"--n-pane-padding-bottom":Q(Z,"bottom"),"--n-font-weight-strong":rt,"--n-tab-color-segment":ot}}),I=w?Ht("tabs",te(()=>`${u.value[0]}${e.type[0]}`),Ce,e):void 0;return Object.assign({mergedClsPrefix:y,mergedValue:B,renderedNames:new Set,segmentCapsuleElRef:D,tabsPaneWrapperRef:J,tabsElRef:C,barElRef:x,addTabInstRef:$,xScrollInstRef:g,scrollWrapperElRef:A,addTabFixed:F,tabWrapperStyle:W,handleNavResize:Ie,mergedSize:u,handleScroll:Ke,handleTabsResize:Ue,cssVars:w?void 0:Ce,themeClass:I==null?void 0:I.themeClass,animationDirection:xe,renderNameListRef:he,yScrollElRef:b,handleSegmentResize:qe,onAnimationBeforeLeave:Me,onAnimationEnter:Ve,onAnimationAfterEnter:He,onRender:I==null?void 0:I.onRender},Ye)},render(){const{mergedClsPrefix:e,type:i,placement:c,addTabFixed:v,addable:l,mergedSize:h,renderNameListRef:y,onRender:w,paneWrapperClass:_,paneWrapperStyle:C,$slots:{default:x,prefix:A,suffix:$}}=this;w==null||w();const g=x?ie(x()).filter(m=>m.type.__TAB_PANE__===!0):[],b=x?ie(x()).filter(m=>m.type.__TAB__===!0):[],s=!b.length,n=i==="card",u=i==="segment",S=!n&&!u&&this.justifyContent;y.value=[];const k=()=>{const m=f("div",{style:this.tabWrapperStyle,class:`${e}-tabs-wrapper`},S?null:f("div",{class:`${e}-tabs-scroll-padding`,style:c==="top"||c==="bottom"?{width:`${this.tabsPadding}px`}:{height:`${this.tabsPadding}px`}}),s?g.map((W,H)=>(y.value.push(W.props.name),be(f(fe,Object.assign({},W.props,{internalCreatedByPane:!0,internalLeftPadded:H!==0&&(!S||S==="center"||S==="start"||S==="end")}),W.children?{default:W.children.tab}:void 0)))):b.map((W,H)=>(y.value.push(W.props.name),be(H!==0&&!S?Le(W):W))),!v&&l&&n?$e(l,(s?g.length:b.length)!==0):null,S?null:f("div",{class:`${e}-tabs-scroll-padding`,style:{width:`${this.tabsPadding}px`}}));return f("div",{ref:"tabsElRef",class:`${e}-tabs-nav-scroll-content`},n&&l?f(se,{onResize:this.handleTabsResize},{default:()=>m}):m,n?f("div",{class:`${e}-tabs-pad`}):null,n?null:f("div",{ref:"barElRef",class:`${e}-tabs-bar`}))},B=u?"top":c;return f("div",{class:[`${e}-tabs`,this.themeClass,`${e}-tabs--${i}-type`,`${e}-tabs--${h}-size`,S&&`${e}-tabs--flex`,`${e}-tabs--${B}`],style:this.cssVars},f("div",{class:[`${e}-tabs-nav--${i}-type`,`${e}-tabs-nav--${B}`,`${e}-tabs-nav`]},we(A,m=>m&&f("div",{class:`${e}-tabs-nav__prefix`},m)),u?f(se,{onResize:this.handleSegmentResize},{default:()=>f("div",{class:`${e}-tabs-rail`,ref:"tabsElRef"},f("div",{class:`${e}-tabs-capsule`,ref:"segmentCapsuleElRef"},f("div",{class:`${e}-tabs-wrapper`},f("div",{class:`${e}-tabs-tab`}))),s?g.map((m,W)=>(y.value.push(m.props.name),f(fe,Object.assign({},m.props,{internalCreatedByPane:!0,internalLeftPadded:W!==0}),m.children?{default:m.children.tab}:void 0))):b.map((m,W)=>(y.value.push(m.props.name),W===0?m:Le(m))))}):f(se,{onResize:this.handleNavResize},{default:()=>f("div",{class:`${e}-tabs-nav-scroll-wrapper`,ref:"scrollWrapperElRef"},["top","bottom"].includes(B)?f(Zt,{ref:"xScrollInstRef",onScroll:this.handleScroll},{default:k}):f("div",{class:`${e}-tabs-nav-y-scroll`,onScroll:this.handleScroll,ref:"yScrollElRef"},k()))}),v&&l&&n?$e(l,!0):null,we($,m=>m&&f("div",{class:`${e}-tabs-nav__suffix`},m))),s&&(this.animated&&(B==="top"||B==="bottom")?f("div",{ref:"tabsPaneWrapperRef",style:C,class:[`${e}-tabs-pane-wrapper`,_]},ze(g,this.mergedValue,this.renderedNames,this.onAnimationBeforeLeave,this.onAnimationEnter,this.onAnimationAfterEnter,this.animationDirection)):ze(g,this.mergedValue,this.renderedNames)))}});function ze(e,i,c,v,l,h,y){const w=[];return e.forEach(_=>{const{name:C,displayDirective:x,"display-directive":A}=_.props,$=b=>x===b||A===b,g=i===C;if(_.key!==void 0&&(_.key=C),g||$("show")||$("show:lazy")&&c.has(C)){c.has(C)||c.add(C);const b=!$("if");w.push(b?Gt(_,[[Dt,g]]):_)}}),y?f(Ft,{name:`${y}-transition`,onBeforeLeave:v,onEnter:l,onAfterEnter:h},{default:()=>w}):w}function $e(e,i){return f(fe,{ref:"addTabInstRef",key:"__addable",name:"__addable",internalCreatedByPane:!0,internalAddable:!0,internalLeftPadded:i,disabled:typeof e=="object"&&e.disabled})}function Le(e){const i=It(e);return i.props?i.props.internalLeftPadded=!0:i.props={internalLeftPadded:!0},i}function be(e){return Array.isArray(e.dynamicProps)?e.dynamicProps.includes("internalLeftPadded")||e.dynamicProps.push("internalLeftPadded"):e.dynamicProps=["internalLeftPadded"],e}const la={class:"monitor-page"},da={class:"monitor-page__header"},ca=["disabled"],ba={class:"monitor-cards"},pa={class:"monitor-card__top"},fa={class:"monitor-card__title"},ua={class:"monitor-card__url"},va={class:"monitor-card__badge"},ga={class:"monitor-card__meta"},ha={key:0,class:"monitor-card__error"},xa={key:1,class:"monitor-card__body"},ma=["disabled","onClick"],ya=Y({__name:"MonitorPage",setup(e){const i=[{id:"onij.fun",label:"onij.fun",cards:[{id:"ping",title:"Ping",url:"https://onij.fun/go/ping"}]},{id:"172.26.250.1",label:"172.26.250.1",cards:[{id:"uvr",title:"UVR",url:"http://172.26.250.1:5555/health"},{id:"whisper",title:"Whisper",url:"http://172.26.250.1:5555/api/v1/whisper/health"}]}],c=L(i[0].id),v=L(!1),l=Nt({});function h(b,s){return`${b}::${s}`}function y(b,s){var n;return((n=l[h(b,s)])==null?void 0:n.status)??"idle"}function w(b,s){var n;return((n=l[h(b,s)])==null?void 0:n.error)??""}function _(b,s){var n;return((n=l[h(b,s)])==null?void 0:n.body)??""}function C(b,s){var u;const n=(u=l[h(b,s)])==null?void 0:u.latencyMs;return typeof n=="number"?`${n} ms`:"—"}function x(b,s){var u;const n=(u=l[h(b,s)])==null?void 0:u.checkedAt;return n?new Date(n).toLocaleTimeString():"—"}function A(b){switch(b){case"loading":return"检测中";case"ok":return"正常";case"fail":return"异常";default:return"未检测"}}async function $(b,s){const n=h(b,s.id);l[n]={status:"loading"};const u=performance.now();try{const S=await fetch(s.url,{method:"GET",cache:"no-store",signal:AbortSignal.timeout(8e3)}),k=Math.round(performance.now()-u),B=(await S.text()).slice(0,500);if(!S.ok){l[n]={status:"fail",latencyMs:k,checkedAt:Date.now(),error:`HTTP ${S.status}`,body:B||void 0};return}l[n]={status:"ok",latencyMs:k,checkedAt:Date.now(),body:B||void 0}}catch(S){let k=S instanceof Error?S.message:String(S);/failed to fetch|networkerror|load failed/i.test(k)&&(k="Failed to fetch（多半未放行）"),l[n]={status:"fail",latencyMs:Math.round(performance.now()-u),checkedAt:Date.now(),error:k}}}async function g(){const b=i.find(s=>s.id===c.value);if(b){v.value=!0;try{await Promise.all(b.cards.map(s=>$(b.id,s)))}finally{v.value=!1}}}return ae(c,()=>{g()}),Ae(()=>{g()}),(b,s)=>(U(),Ut(yt,null,{default:de(()=>[R("div",la,[R("div",da,[s[1]||(s[1]=R("h1",{class:"browse-section-title border-l-4 border-blue-500 pl-2"},"监控",-1)),R("button",{type:"button",class:"monitor-btn",disabled:v.value,onClick:g},G(v.value?"检测中…":"刷新当前标签"),9,ca)]),Pe(Re(sa),{value:c.value,"onUpdate:value":s[0]||(s[0]=n=>c.value=n),type:"line",animated:"",class:"monitor-tabs"},{default:de(()=>[(U(),K(pe,null,_e(i,n=>Pe(Re(ra),{key:n.id,name:n.id,tab:n.label},{default:de(()=>[R("div",ba,[(U(!0),K(pe,null,_e(n.cards,u=>(U(),K("article",{key:u.id,class:Xt(["monitor-card",`monitor-card--${y(n.id,u.id)}`])},[R("div",pa,[R("div",null,[R("h2",fa,G(u.title),1),R("p",ua,G(u.url),1)]),R("span",va,G(A(y(n.id,u.id))),1)]),R("dl",ga,[R("div",null,[s[2]||(s[2]=R("dt",null,"延迟",-1)),R("dd",null,G(C(n.id,u.id)),1)]),R("div",null,[s[3]||(s[3]=R("dt",null,"最近检测",-1)),R("dd",null,G(x(n.id,u.id)),1)])]),w(n.id,u.id)?(U(),K("p",ha,G(w(n.id,u.id)),1)):_(n.id,u.id)?(U(),K("pre",xa,G(_(n.id,u.id)),1)):Kt("",!0),R("button",{type:"button",class:"monitor-btn monitor-btn--ghost",disabled:v.value,onClick:S=>$(n.id,u)}," 单独检测 ",8,ma)],2))),128))])]),_:2},1032,["name","tab"])),64))]),_:1},8,["value"])])]),_:1}))}}),Ta=Yt(ya,[["__scopeId","data-v-9e4eaf56"]]);export{Ta as default};
