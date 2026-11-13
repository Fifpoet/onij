import{d as re,aa as g,r as z,bK as $e,bL as Rt,bM as kt,bN as zt,bO as $t,ab as Lt,av as Bt,a_ as At,a5 as Ve,aN as Wt,b3 as Mt,F as de,aF as Et,ar as jt,c as ce,az as Ot,ad as o,ah as c,ac as k,ai as A,aD as Vt,bl as ve,aj as Le,bP as ge,ak as Ht,al as He,bx as Be,aG as Gt,w as be,o as Ge,aC as Dt,aw as N,bQ as It,aH as Ft,an as U,ao as ie,ap as Nt,bo as Ut,A as Xt,bc as Kt,bR as Yt,ba as qt,ax as he,as as le,bS as De,bT as Jt,bU as Qt,a as Zt,e as ea,f as se,i as b,n as q,t as W,g as X,j as me,y as E,m as ta,l as I,I as Ae,h as xe,_ as aa}from"./index-DC-8R5Yo.js";import{P as ra}from"./PageContent-BrMQ_Y4m.js";import{A as oa,d as ye,f as na}from"./Select-Bhx6Ao56.js";const ia=$e(".v-x-scroll",{overflow:"auto",scrollbarWidth:"none"},[$e("&::-webkit-scrollbar",{width:0,height:0})]),la=re({name:"XScroll",props:{disabled:Boolean,onScroll:Function},setup(){const e=z(null);function n(f){!(f.currentTarget.offsetWidth<f.currentTarget.scrollWidth)||f.deltaY===0||(f.currentTarget.scrollLeft+=f.deltaY+f.deltaX,f.preventDefault())}const s=Rt();return ia.mount({id:"vueuc/x-scroll",head:!0,anchorMetaName:kt,ssr:s}),Object.assign({selfRef:e,handleWheel:n},{scrollTo(...f){var h;(h=e.value)===null||h===void 0||h.scrollTo(...f)}})},render(){return g("div",{ref:"selfRef",onScroll:this.onScroll,onWheel:this.disabled?void 0:this.handleWheel,class:"v-x-scroll"},this.$slots)}});var sa="Expected a function";function Se(e,n,s){var p=!0,f=!0;if(typeof e!="function")throw new TypeError(sa);return $t(s)&&(p="leading"in s?!!s.leading:p,f="trailing"in s?!!s.trailing:f),zt(e,n,{leading:p,maxWait:n,trailing:f})}const da={tabFontSizeSmall:"14px",tabFontSizeMedium:"14px",tabFontSizeLarge:"16px",tabGapSmallLine:"36px",tabGapMediumLine:"36px",tabGapLargeLine:"36px",tabGapSmallLineVertical:"8px",tabGapMediumLineVertical:"8px",tabGapLargeLineVertical:"8px",tabPaddingSmallLine:"6px 0",tabPaddingMediumLine:"10px 0",tabPaddingLargeLine:"14px 0",tabPaddingVerticalSmallLine:"6px 12px",tabPaddingVerticalMediumLine:"8px 16px",tabPaddingVerticalLargeLine:"10px 20px",tabGapSmallBar:"36px",tabGapMediumBar:"36px",tabGapLargeBar:"36px",tabGapSmallBarVertical:"8px",tabGapMediumBarVertical:"8px",tabGapLargeBarVertical:"8px",tabPaddingSmallBar:"4px 0",tabPaddingMediumBar:"6px 0",tabPaddingLargeBar:"10px 0",tabPaddingVerticalSmallBar:"6px 12px",tabPaddingVerticalMediumBar:"8px 16px",tabPaddingVerticalLargeBar:"10px 20px",tabGapSmallCard:"4px",tabGapMediumCard:"4px",tabGapLargeCard:"4px",tabGapSmallCardVertical:"4px",tabGapMediumCardVertical:"4px",tabGapLargeCardVertical:"4px",tabPaddingSmallCard:"8px 16px",tabPaddingMediumCard:"10px 20px",tabPaddingLargeCard:"12px 24px",tabPaddingSmallSegment:"4px 0",tabPaddingMediumSegment:"6px 0",tabPaddingLargeSegment:"8px 0",tabPaddingVerticalLargeSegment:"0 8px",tabPaddingVerticalSmallCard:"8px 12px",tabPaddingVerticalMediumCard:"10px 16px",tabPaddingVerticalLargeCard:"12px 20px",tabPaddingVerticalSmallSegment:"0 4px",tabPaddingVerticalMediumSegment:"0 6px",tabGapSmallSegment:"0",tabGapMediumSegment:"0",tabGapLargeSegment:"0",tabGapSmallSegmentVertical:"0",tabGapMediumSegmentVertical:"0",tabGapLargeSegmentVertical:"0",panePaddingSmall:"8px 0 0 0",panePaddingMedium:"12px 0 0 0",panePaddingLarge:"16px 0 0 0",closeSize:"18px",closeIconSize:"14px"};function ca(e){const{textColor2:n,primaryColor:s,textColorDisabled:p,closeIconColor:f,closeIconColorHover:h,closeIconColorPressed:u,closeColorHover:C,closeColorPressed:_,tabColor:x,baseColor:m,dividerColor:j,fontWeight:$,textColor1:y,borderRadius:P,fontSize:T,fontWeightStrong:B}=e;return Object.assign(Object.assign({},da),{colorSegment:x,tabFontSizeCard:T,tabTextColorLine:y,tabTextColorActiveLine:s,tabTextColorHoverLine:s,tabTextColorDisabledLine:p,tabTextColorSegment:y,tabTextColorActiveSegment:n,tabTextColorHoverSegment:n,tabTextColorDisabledSegment:p,tabTextColorBar:y,tabTextColorActiveBar:s,tabTextColorHoverBar:s,tabTextColorDisabledBar:p,tabTextColorCard:y,tabTextColorHoverCard:y,tabTextColorActiveCard:s,tabTextColorDisabledCard:p,barColor:s,closeIconColor:f,closeIconColorHover:h,closeIconColorPressed:u,closeColorHover:C,closeColorPressed:_,closeBorderRadius:P,tabColor:x,tabColorSegment:m,tabBorderColor:j,tabFontWeightActive:$,tabFontWeight:$,tabBorderRadius:P,paneTextColor:n,fontWeightStrong:B})}const ba={common:Lt,self:ca},_e=Bt("n-tabs"),Ie={tab:[String,Number,Object,Function],name:{type:[String,Number],required:!0},disabled:Boolean,displayDirective:{type:String,default:"if"},closable:{type:Boolean,default:void 0},tabProps:Object,label:[String,Number,Object,Function]},We=re({__TAB_PANE__:!0,name:"TabPane",alias:["TabPanel"],props:Ie,slots:Object,setup(e){const n=Ve(_e,null);return n||At("tab-pane","`n-tab-pane` must be placed inside `n-tabs`."),{style:n.paneStyleRef,class:n.paneClassRef,mergedClsPrefix:n.mergedClsPrefixRef}},render(){return g("div",{class:[`${this.mergedClsPrefix}-tab-pane`,this.class],style:this.style},this.$slots)}}),pa=Object.assign({internalLeftPadded:Boolean,internalAddable:Boolean,internalCreatedByPane:Boolean},Ot(Ie,["displayDirective"])),we=re({__TAB__:!0,inheritAttrs:!1,name:"Tab",props:pa,setup(e){const{mergedClsPrefixRef:n,valueRef:s,typeRef:p,closableRef:f,tabStyleRef:h,addTabStyleRef:u,tabClassRef:C,addTabClassRef:_,tabChangeIdRef:x,onBeforeLeaveRef:m,triggerRef:j,handleAdd:$,activateTab:y,handleClose:P}=Ve(_e);return{trigger:j,mergedClosable:ce(()=>{if(e.internalAddable)return!1;const{closable:T}=e;return T===void 0?f.value:T}),style:h,addStyle:u,tabClass:C,addTabClass:_,clsPrefix:n,value:s,type:p,handleClose(T){T.stopPropagation(),!e.disabled&&P(e.name)},activateTab(){if(e.disabled)return;if(e.internalAddable){$();return}const{name:T}=e,B=++x.id;if(T!==s.value){const{value:V}=m;V?Promise.resolve(V(e.name,s.value)).then(M=>{M&&x.id===B&&y(T)}):y(T)}}}},render(){const{internalAddable:e,clsPrefix:n,name:s,disabled:p,label:f,tab:h,value:u,mergedClosable:C,trigger:_,$slots:{default:x}}=this,m=f??h;return g("div",{class:`${n}-tabs-tab-wrapper`},this.internalLeftPadded?g("div",{class:`${n}-tabs-tab-pad`}):null,g("div",Object.assign({key:s,"data-name":s,"data-disabled":p?!0:void 0},Wt({class:[`${n}-tabs-tab`,u===s&&`${n}-tabs-tab--active`,p&&`${n}-tabs-tab--disabled`,C&&`${n}-tabs-tab--closable`,e&&`${n}-tabs-tab--addable`,e?this.addTabClass:this.tabClass],onClick:_==="click"?this.activateTab:void 0,onMouseenter:_==="hover"?this.activateTab:void 0,style:e?this.addStyle:this.style},this.internalCreatedByPane?this.tabProps||{}:this.$attrs)),g("span",{class:`${n}-tabs-tab__label`},e?g(de,null,g("div",{class:`${n}-tabs-tab__height-placeholder`}," "),g(Et,{clsPrefix:n},{default:()=>g(oa,null)})):x?x():typeof m=="object"?m:Mt(m??s)),C&&this.type==="card"?g(jt,{clsPrefix:n,class:`${n}-tabs-tab__close`,onClick:this.handleClose,disabled:p}):null))}}),fa=o("tabs",`
 box-sizing: border-box;
 width: 100%;
 display: flex;
 flex-direction: column;
 transition:
 background-color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
`,[c("segment-type",[o("tabs-rail",[k("&.transition-disabled",[o("tabs-capsule",`
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
 `),k("&:hover",`
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
 `,[A("prefix, suffix",`
 display: flex;
 align-items: center;
 `),A("prefix","padding-right: 16px;"),A("suffix","padding-left: 16px;")]),c("top, bottom",[o("tabs-nav-scroll-wrapper",[k("&::before",`
 top: 0;
 bottom: 0;
 left: 0;
 width: 20px;
 `),k("&::after",`
 top: 0;
 bottom: 0;
 right: 0;
 width: 20px;
 `),c("shadow-start",[k("&::before",`
 box-shadow: inset 10px 0 8px -8px rgba(0, 0, 0, .12);
 `)]),c("shadow-end",[k("&::after",`
 box-shadow: inset -10px 0 8px -8px rgba(0, 0, 0, .12);
 `)])])]),c("left, right",[o("tabs-nav-scroll-content",`
 flex-direction: column;
 `),o("tabs-nav-scroll-wrapper",[k("&::before",`
 top: 0;
 left: 0;
 right: 0;
 height: 20px;
 `),k("&::after",`
 bottom: 0;
 left: 0;
 right: 0;
 height: 20px;
 `),c("shadow-start",[k("&::before",`
 box-shadow: inset 0 10px 8px -8px rgba(0, 0, 0, .12);
 `)]),c("shadow-end",[k("&::after",`
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
 `,[k("&::-webkit-scrollbar, &::-webkit-scrollbar-track-piece, &::-webkit-scrollbar-thumb",`
 width: 0;
 height: 0;
 display: none;
 `)]),k("&::before, &::after",`
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
 `,[c("disabled",{cursor:"not-allowed"}),A("close",`
 margin-left: 6px;
 transition:
 background-color .3s var(--n-bezier),
 color .3s var(--n-bezier);
 `),A("label",`
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
 `,[k("&.transition-disabled",`
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
 `,[k("&.next-transition-leave-active, &.prev-transition-leave-active, &.next-transition-enter-active, &.prev-transition-enter-active",`
 transition:
 color .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 transform .2s var(--n-bezier),
 opacity .2s var(--n-bezier);
 `),k("&.next-transition-leave-active, &.prev-transition-leave-active",`
 position: absolute;
 `),k("&.next-transition-enter-from, &.prev-transition-leave-to",`
 transform: translateX(32px);
 opacity: 0;
 `),k("&.next-transition-leave-to, &.prev-transition-enter-from",`
 transform: translateX(-32px);
 opacity: 0;
 `),k("&.next-transition-leave-from, &.next-transition-enter-to, &.prev-transition-leave-from, &.prev-transition-enter-to",`
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
 `,[k("&:hover",{color:"var(--n-tab-text-color-hover)"}),c("active",`
 color: var(--n-tab-text-color-active);
 font-weight: var(--n-tab-font-weight-active);
 `),c("disabled",{color:"var(--n-tab-text-color-disabled)"})])]),o("tabs-nav",[c("line-type",[c("top",[A("prefix, suffix",`
 border-bottom: 1px solid var(--n-tab-border-color);
 `),o("tabs-nav-scroll-content",`
 border-bottom: 1px solid var(--n-tab-border-color);
 `),o("tabs-bar",`
 bottom: -1px;
 `)]),c("left",[A("prefix, suffix",`
 border-right: 1px solid var(--n-tab-border-color);
 `),o("tabs-nav-scroll-content",`
 border-right: 1px solid var(--n-tab-border-color);
 `),o("tabs-bar",`
 right: -1px;
 `)]),c("right",[A("prefix, suffix",`
 border-left: 1px solid var(--n-tab-border-color);
 `),o("tabs-nav-scroll-content",`
 border-left: 1px solid var(--n-tab-border-color);
 `),o("tabs-bar",`
 left: -1px;
 `)]),c("bottom",[A("prefix, suffix",`
 border-top: 1px solid var(--n-tab-border-color);
 `),o("tabs-nav-scroll-content",`
 border-top: 1px solid var(--n-tab-border-color);
 `),o("tabs-bar",`
 top: -1px;
 `)]),A("prefix, suffix",`
 transition: border-color .3s var(--n-bezier);
 `),o("tabs-nav-scroll-content",`
 transition: border-color .3s var(--n-bezier);
 `),o("tabs-bar",`
 border-radius: 0;
 `)]),c("card-type",[A("prefix, suffix",`
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
 `,[A("height-placeholder",`
 width: 0;
 font-size: var(--n-tab-font-size);
 `),Vt("disabled",[k("&:hover",`
 color: var(--n-tab-text-color-hover);
 `)])]),c("closable","padding-right: 8px;"),c("active",`
 background-color: #0000;
 font-weight: var(--n-tab-font-weight-active);
 color: var(--n-tab-text-color-active);
 `),c("disabled","color: var(--n-tab-text-color-disabled);")])]),c("left, right",`
 flex-direction: column; 
 `,[A("prefix, suffix",`
 padding: var(--n-tab-padding-vertical);
 `),o("tabs-wrapper",`
 flex-direction: column;
 `),o("tabs-tab-wrapper",`
 flex-direction: column;
 `,[o("tabs-tab-pad",`
 height: var(--n-tab-gap-vertical);
 width: 100%;
 `)])]),c("top",[c("card-type",[o("tabs-scroll-padding","border-bottom: 1px solid var(--n-tab-border-color);"),A("prefix, suffix",`
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
 `)])]),c("left",[c("card-type",[o("tabs-scroll-padding","border-right: 1px solid var(--n-tab-border-color);"),A("prefix, suffix",`
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
 `)])]),c("right",[c("card-type",[o("tabs-scroll-padding","border-left: 1px solid var(--n-tab-border-color);"),A("prefix, suffix",`
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
 `)])]),c("bottom",[c("card-type",[o("tabs-scroll-padding","border-top: 1px solid var(--n-tab-border-color);"),A("prefix, suffix",`
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
 `)])])])]),ua=Object.assign(Object.assign({},He.props),{value:[String,Number],defaultValue:[String,Number],trigger:{type:String,default:"click"},type:{type:String,default:"bar"},closable:Boolean,justifyContent:String,size:{type:String,default:"medium"},placement:{type:String,default:"top"},tabStyle:[String,Object],tabClass:String,addTabStyle:[String,Object],addTabClass:String,barWidth:Number,paneClass:String,paneStyle:[String,Object],paneWrapperClass:String,paneWrapperStyle:[String,Object],addable:[Boolean,Object],tabsPadding:{type:Number,default:0},animated:Boolean,onBeforeLeave:Function,onAdd:Function,"onUpdate:value":[Function,Array],onUpdateValue:[Function,Array],onClose:[Function,Array],labelSize:String,activeName:[String,Number],onActiveNameChange:[Function,Array]}),va=re({name:"Tabs",props:ua,slots:Object,setup(e,{slots:n}){var s,p,f,h;const{mergedClsPrefixRef:u,inlineThemeDisabled:C}=Ht(e),_=He("Tabs","-tabs",fa,ba,e,u),x=z(null),m=z(null),j=z(null),$=z(null),y=z(null),P=z(null),T=z(!0),B=z(!0),V=Be(e,["labelSize","size"]),M=Be(e,["activeName","value"]),Y=z((p=(s=M.value)!==null&&s!==void 0?s:e.defaultValue)!==null&&p!==void 0?p:n.default?(h=(f=ve(n.default())[0])===null||f===void 0?void 0:f.props)===null||h===void 0?void 0:h.name:null),O=Gt(M,Y),S={id:0},d=ce(()=>{if(!(!e.justifyContent||e.type==="card"))return{display:"flex",justifyContent:e.justifyContent}});be(O,()=>{S.id=0,H(),te()});function i(){var t;const{value:a}=O;return a===null?null:(t=x.value)===null||t===void 0?void 0:t.querySelector(`[data-name="${a}"]`)}function r(t){if(e.type==="card")return;const{value:a}=m;if(!a)return;const l=a.style.opacity==="0";if(t){const v=`${u.value}-tabs-bar--disabled`,{barWidth:R,placement:G}=e;if(t.dataset.disabled==="true"?a.classList.add(v):a.classList.remove(v),["top","bottom"].includes(G)){if(L(["top","maxHeight","height"]),typeof R=="number"&&t.offsetWidth>=R){const D=Math.floor((t.offsetWidth-R)/2)+t.offsetLeft;a.style.left=`${D}px`,a.style.maxWidth=`${R}px`}else a.style.left=`${t.offsetLeft}px`,a.style.maxWidth=`${t.offsetWidth}px`;a.style.width="8192px",l&&(a.style.transition="none"),a.offsetWidth,l&&(a.style.transition="",a.style.opacity="1")}else{if(L(["left","maxWidth","width"]),typeof R=="number"&&t.offsetHeight>=R){const D=Math.floor((t.offsetHeight-R)/2)+t.offsetTop;a.style.top=`${D}px`,a.style.maxHeight=`${R}px`}else a.style.top=`${t.offsetTop}px`,a.style.maxHeight=`${t.offsetHeight}px`;a.style.height="8192px",l&&(a.style.transition="none"),a.offsetHeight,l&&(a.style.transition="",a.style.opacity="1")}}}function w(){if(e.type==="card")return;const{value:t}=m;t&&(t.style.opacity="0")}function L(t){const{value:a}=m;if(a)for(const l of t)a.style[l]=""}function H(){if(e.type==="card")return;const t=i();t?r(t):w()}function te(){var t;const a=(t=y.value)===null||t===void 0?void 0:t.$el;if(!a)return;const l=i();if(!l)return;const{scrollLeft:v,offsetWidth:R}=a,{offsetLeft:G,offsetWidth:D}=l;v>G?a.scrollTo({top:0,left:G,behavior:"smooth"}):G+D>v+R&&a.scrollTo({top:0,left:G+D-R,behavior:"smooth"})}const oe=z(null);let pe=0,K=null;function Fe(t){const a=oe.value;if(a){pe=t.getBoundingClientRect().height;const l=`${pe}px`,v=()=>{a.style.height=l,a.style.maxHeight=l};K?(v(),K(),K=null):K=v}}function Ne(t){const a=oe.value;if(a){const l=t.getBoundingClientRect().height,v=()=>{document.body.offsetHeight,a.style.maxHeight=`${l}px`,a.style.height=`${Math.max(pe,l)}px`};K?(K(),K=null,v()):K=v}}function Ue(){const t=oe.value;if(t){t.style.maxHeight="",t.style.height="";const{paneWrapperStyle:a}=e;if(typeof a=="string")t.style.cssText=a;else if(a){const{maxHeight:l,height:v}=a;l!==void 0&&(t.style.maxHeight=l),v!==void 0&&(t.style.height=v)}}}const Te={value:[]},Pe=z("next");function Xe(t){const a=O.value;let l="next";for(const v of Te.value){if(v===a)break;if(v===t){l="prev";break}}Pe.value=l,Ke(t)}function Ke(t){const{onActiveNameChange:a,onUpdateValue:l,"onUpdate:value":v}=e;a&&le(a,t),l&&le(l,t),v&&le(v,t),Y.value=t}function Ye(t){const{onClose:a}=e;a&&le(a,t)}function Re(){const{value:t}=m;if(!t)return;const a="transition-disabled";t.classList.add(a),H(),t.classList.remove(a)}const J=z(null);function fe({transitionDisabled:t}){const a=x.value;if(!a)return;t&&a.classList.add("transition-disabled");const l=i();l&&J.value&&(J.value.style.width=`${l.offsetWidth}px`,J.value.style.height=`${l.offsetHeight}px`,J.value.style.transform=`translateX(${l.offsetLeft-Ut(getComputedStyle(a).paddingLeft)}px)`,t&&J.value.offsetWidth),t&&a.classList.remove("transition-disabled")}be([O],()=>{e.type==="segment"&&he(()=>{fe({transitionDisabled:!1})})}),Ge(()=>{e.type==="segment"&&fe({transitionDisabled:!0})});let ke=0;function qe(t){var a;if(t.contentRect.width===0&&t.contentRect.height===0||ke===t.contentRect.width)return;ke=t.contentRect.width;const{type:l}=e;if((l==="line"||l==="bar")&&Re(),l!=="segment"){const{placement:v}=e;ue((v==="top"||v==="bottom"?(a=y.value)===null||a===void 0?void 0:a.$el:P.value)||null)}}const Je=Se(qe,64);be([()=>e.justifyContent,()=>e.size],()=>{he(()=>{const{type:t}=e;(t==="line"||t==="bar")&&Re()})});const Q=z(!1);function Qe(t){var a;const{target:l,contentRect:{width:v,height:R}}=t,G=l.parentElement.parentElement.offsetWidth,D=l.parentElement.parentElement.offsetHeight,{placement:ee}=e;if(!Q.value)ee==="top"||ee==="bottom"?G<v&&(Q.value=!0):D<R&&(Q.value=!0);else{const{value:ae}=$;if(!ae)return;ee==="top"||ee==="bottom"?G-v>ae.$el.offsetWidth&&(Q.value=!1):D-R>ae.$el.offsetHeight&&(Q.value=!1)}ue(((a=y.value)===null||a===void 0?void 0:a.$el)||null)}const Ze=Se(Qe,64);function et(){const{onAdd:t}=e;t&&t(),he(()=>{const a=i(),{value:l}=y;!a||!l||l.scrollTo({left:a.offsetLeft,top:0,behavior:"smooth"})})}function ue(t){if(!t)return;const{placement:a}=e;if(a==="top"||a==="bottom"){const{scrollLeft:l,scrollWidth:v,offsetWidth:R}=t;T.value=l<=0,B.value=l+R>=v}else{const{scrollTop:l,scrollHeight:v,offsetHeight:R}=t;T.value=l<=0,B.value=l+R>=v}}const tt=Se(t=>{ue(t.target)},64);Dt(_e,{triggerRef:N(e,"trigger"),tabStyleRef:N(e,"tabStyle"),tabClassRef:N(e,"tabClass"),addTabStyleRef:N(e,"addTabStyle"),addTabClassRef:N(e,"addTabClass"),paneClassRef:N(e,"paneClass"),paneStyleRef:N(e,"paneStyle"),mergedClsPrefixRef:u,typeRef:N(e,"type"),closableRef:N(e,"closable"),valueRef:O,tabChangeIdRef:S,onBeforeLeaveRef:N(e,"onBeforeLeave"),activateTab:Xe,handleClose:Ye,handleAdd:et}),It(()=>{H(),te()}),Ft(()=>{const{value:t}=j;if(!t)return;const{value:a}=u,l=`${a}-tabs-nav-scroll-wrapper--shadow-start`,v=`${a}-tabs-nav-scroll-wrapper--shadow-end`;T.value?t.classList.remove(l):t.classList.add(l),B.value?t.classList.remove(v):t.classList.add(v)});const at={syncBarPosition:()=>{H()}},rt=()=>{fe({transitionDisabled:!0})},ze=ce(()=>{const{value:t}=V,{type:a}=e,l={card:"Card",bar:"Bar",line:"Line",segment:"Segment"}[a],v=`${t}${l}`,{self:{barColor:R,closeIconColor:G,closeIconColorHover:D,closeIconColorPressed:ee,tabColor:ae,tabBorderColor:ot,paneTextColor:nt,tabFontWeight:it,tabBorderRadius:lt,tabFontWeightActive:st,colorSegment:dt,fontWeightStrong:ct,tabColorSegment:bt,closeSize:pt,closeIconSize:ft,closeColorHover:ut,closeColorPressed:vt,closeBorderRadius:gt,[U("panePadding",t)]:ne,[U("tabPadding",v)]:ht,[U("tabPaddingVertical",v)]:mt,[U("tabGap",v)]:xt,[U("tabGap",`${v}Vertical`)]:yt,[U("tabTextColor",a)]:St,[U("tabTextColorActive",a)]:Ct,[U("tabTextColorHover",a)]:wt,[U("tabTextColorDisabled",a)]:_t,[U("tabFontSize",t)]:Tt},common:{cubicBezierEaseInOut:Pt}}=_.value;return{"--n-bezier":Pt,"--n-color-segment":dt,"--n-bar-color":R,"--n-tab-font-size":Tt,"--n-tab-text-color":St,"--n-tab-text-color-active":Ct,"--n-tab-text-color-disabled":_t,"--n-tab-text-color-hover":wt,"--n-pane-text-color":nt,"--n-tab-border-color":ot,"--n-tab-border-radius":lt,"--n-close-size":pt,"--n-close-icon-size":ft,"--n-close-color-hover":ut,"--n-close-color-pressed":vt,"--n-close-border-radius":gt,"--n-close-icon-color":G,"--n-close-icon-color-hover":D,"--n-close-icon-color-pressed":ee,"--n-tab-color":ae,"--n-tab-font-weight":it,"--n-tab-font-weight-active":st,"--n-tab-padding":ht,"--n-tab-padding-vertical":mt,"--n-tab-gap":xt,"--n-tab-gap-vertical":yt,"--n-pane-padding-left":ie(ne,"left"),"--n-pane-padding-right":ie(ne,"right"),"--n-pane-padding-top":ie(ne,"top"),"--n-pane-padding-bottom":ie(ne,"bottom"),"--n-font-weight-strong":ct,"--n-tab-color-segment":bt}}),Z=C?Nt("tabs",ce(()=>`${V.value[0]}${e.type[0]}`),ze,e):void 0;return Object.assign({mergedClsPrefix:u,mergedValue:O,renderedNames:new Set,segmentCapsuleElRef:J,tabsPaneWrapperRef:oe,tabsElRef:x,barElRef:m,addTabInstRef:$,xScrollInstRef:y,scrollWrapperElRef:j,addTabFixed:Q,tabWrapperStyle:d,handleNavResize:Je,mergedSize:V,handleScroll:tt,handleTabsResize:Ze,cssVars:C?void 0:ze,themeClass:Z==null?void 0:Z.themeClass,animationDirection:Pe,renderNameListRef:Te,yScrollElRef:P,handleSegmentResize:rt,onAnimationBeforeLeave:Fe,onAnimationEnter:Ne,onAnimationAfterEnter:Ue,onRender:Z==null?void 0:Z.onRender},at)},render(){const{mergedClsPrefix:e,type:n,placement:s,addTabFixed:p,addable:f,mergedSize:h,renderNameListRef:u,onRender:C,paneWrapperClass:_,paneWrapperStyle:x,$slots:{default:m,prefix:j,suffix:$}}=this;C==null||C();const y=m?ve(m()).filter(S=>S.type.__TAB_PANE__===!0):[],P=m?ve(m()).filter(S=>S.type.__TAB__===!0):[],T=!P.length,B=n==="card",V=n==="segment",M=!B&&!V&&this.justifyContent;u.value=[];const Y=()=>{const S=g("div",{style:this.tabWrapperStyle,class:`${e}-tabs-wrapper`},M?null:g("div",{class:`${e}-tabs-scroll-padding`,style:s==="top"||s==="bottom"?{width:`${this.tabsPadding}px`}:{height:`${this.tabsPadding}px`}}),T?y.map((d,i)=>(u.value.push(d.props.name),Ce(g(we,Object.assign({},d.props,{internalCreatedByPane:!0,internalLeftPadded:i!==0&&(!M||M==="center"||M==="start"||M==="end")}),d.children?{default:d.children.tab}:void 0)))):P.map((d,i)=>(u.value.push(d.props.name),Ce(i!==0&&!M?je(d):d))),!p&&f&&B?Ee(f,(T?y.length:P.length)!==0):null,M?null:g("div",{class:`${e}-tabs-scroll-padding`,style:{width:`${this.tabsPadding}px`}}));return g("div",{ref:"tabsElRef",class:`${e}-tabs-nav-scroll-content`},B&&f?g(ge,{onResize:this.handleTabsResize},{default:()=>S}):S,B?g("div",{class:`${e}-tabs-pad`}):null,B?null:g("div",{ref:"barElRef",class:`${e}-tabs-bar`}))},O=V?"top":s;return g("div",{class:[`${e}-tabs`,this.themeClass,`${e}-tabs--${n}-type`,`${e}-tabs--${h}-size`,M&&`${e}-tabs--flex`,`${e}-tabs--${O}`],style:this.cssVars},g("div",{class:[`${e}-tabs-nav--${n}-type`,`${e}-tabs-nav--${O}`,`${e}-tabs-nav`]},Le(j,S=>S&&g("div",{class:`${e}-tabs-nav__prefix`},S)),V?g(ge,{onResize:this.handleSegmentResize},{default:()=>g("div",{class:`${e}-tabs-rail`,ref:"tabsElRef"},g("div",{class:`${e}-tabs-capsule`,ref:"segmentCapsuleElRef"},g("div",{class:`${e}-tabs-wrapper`},g("div",{class:`${e}-tabs-tab`}))),T?y.map((S,d)=>(u.value.push(S.props.name),g(we,Object.assign({},S.props,{internalCreatedByPane:!0,internalLeftPadded:d!==0}),S.children?{default:S.children.tab}:void 0))):P.map((S,d)=>(u.value.push(S.props.name),d===0?S:je(S))))}):g(ge,{onResize:this.handleNavResize},{default:()=>g("div",{class:`${e}-tabs-nav-scroll-wrapper`,ref:"scrollWrapperElRef"},["top","bottom"].includes(O)?g(la,{ref:"xScrollInstRef",onScroll:this.handleScroll},{default:Y}):g("div",{class:`${e}-tabs-nav-y-scroll`,onScroll:this.handleScroll,ref:"yScrollElRef"},Y()))}),p&&f&&B?Ee(f,!0):null,Le($,S=>S&&g("div",{class:`${e}-tabs-nav__suffix`},S))),T&&(this.animated&&(O==="top"||O==="bottom")?g("div",{ref:"tabsPaneWrapperRef",style:x,class:[`${e}-tabs-pane-wrapper`,_]},Me(y,this.mergedValue,this.renderedNames,this.onAnimationBeforeLeave,this.onAnimationEnter,this.onAnimationAfterEnter,this.animationDirection)):Me(y,this.mergedValue,this.renderedNames)))}});function Me(e,n,s,p,f,h,u){const C=[];return e.forEach(_=>{const{name:x,displayDirective:m,"display-directive":j}=_.props,$=P=>m===P||j===P,y=n===x;if(_.key!==void 0&&(_.key=x),y||$("show")||$("show:lazy")&&s.has(x)){s.has(x)||s.add(x);const P=!$("if");C.push(P?Xt(_,[[Kt,y]]):_)}}),u?g(Yt,{name:`${u}-transition`,onBeforeLeave:p,onEnter:f,onAfterEnter:h},{default:()=>C}):C}function Ee(e,n){return g(we,{ref:"addTabInstRef",key:"__addable",name:"__addable",internalCreatedByPane:!0,internalAddable:!0,internalLeftPadded:n,disabled:typeof e=="object"&&e.disabled})}function je(e){const n=qt(e);return n.props?n.props.internalLeftPadded=!0:n.props={internalLeftPadded:!0},n}function Ce(e){return Array.isArray(e.dynamicProps)?e.dynamicProps.includes("internalLeftPadded")||e.dynamicProps.push("internalLeftPadded"):e.dynamicProps=["internalLeftPadded"],e}const ga=8e3,ha="185811";function ma(e){const n=De(e);return e.kind==="gdstudio"?`${n}/time`:`${n}/?server=netease&type=url&id=${ha}`}function Oe(e){return De(e)||Jt(e.kind)}function xa(e){return typeof e=="string"&&/^https?:\/\//i.test(e)}function ya(e,n){if(e<200||e>=400)return{ok:!1,error:`HTTP ${e}`};const s=n.trim();if(!s)return{ok:!1,error:"空响应"};if(/^https?:\/\//i.test(s))return{ok:!0};try{const p=JSON.parse(s);return xa(p==null?void 0:p.url)?{ok:!0}:{ok:!1,error:"无播放地址"}}catch{return/<html/i.test(s)?{ok:!1,error:"返回了网页"}:{ok:!1,error:"无法解析"}}}function Sa(e,n){return e<200||e>=400?{ok:!1,error:`HTTP ${e}`}:/^\d+$/.test(n.trim())?{ok:!0}:{ok:!1,error:"time 无效"}}async function Ca(e){const n=ma(e),s=performance.now(),p=new AbortController,f=setTimeout(()=>p.abort(),ga);try{const h=await fetch(n,{method:"GET",cache:"no-store",redirect:e.kind==="meting"?"manual":"follow",signal:p.signal}),u=Math.round(performance.now()-s),C=h.headers.get("Location");if(e.kind==="meting"&&C&&/^https?:\/\//i.test(C))return{ok:!0,latencyMs:u,body:C.slice(0,500)};const _=(await h.text()).slice(0,500),x=e.kind==="gdstudio"?Sa(h.status,_):ya(h.status,_);return x.ok?{ok:!0,latencyMs:u,body:_||void 0}:{ok:!1,latencyMs:u,error:x.error,body:_||void 0}}catch(h){let u=h instanceof Error?h.message:String(h);return/failed to fetch|networkerror|load failed/i.test(u)&&(u="Failed to fetch（多半未放行或源已挂）"),{ok:!1,latencyMs:Math.round(performance.now()-s),error:u}}finally{clearTimeout(f)}}const wa={class:"monitor-page"},_a={class:"monitor-page__header"},Ta=["disabled"],Pa={class:"monitor-cards"},Ra={class:"monitor-card__top"},ka={class:"monitor-card__title"},za={class:"monitor-card__url"},$a={class:"monitor-card__badge"},La={class:"monitor-card__meta"},Ba={key:0,class:"monitor-card__error"},Aa={key:1,class:"monitor-card__body"},Wa=["disabled","onClick"],Ma={class:"monitor-cards"},Ea={class:"monitor-card__top"},ja={class:"monitor-card__title"},Oa={class:"monitor-card__url"},Va={class:"monitor-card__badge"},Ha={class:"monitor-field"},Ga={class:"monitor-card__meta"},Da={key:0,class:"monitor-card__error"},Ia={key:1,class:"monitor-card__body"},Fa={class:"monitor-card__actions"},Na=["disabled","onClick"],Ua=["onClick"],Xa=["disabled","onClick"],Ka=["disabled","onClick"],Ya=["onClick"],qa=["disabled"],F="play-sources",Ja=re({__name:"MonitorPage",setup(e){const n=[{id:"onij.fun",label:"onij.fun",cards:[{id:"ping",title:"Ping",url:"https://onij.fun/go/ping"}]},{id:"172.26.250.1",label:"172.26.250.1",cards:[{id:"uvr",title:"UVR",url:"http://172.26.250.1:5555/health"},{id:"whisper",title:"Whisper",url:"http://172.26.250.1:5555/api/v1/whisper/health"}]}],s=[{label:"Meting",value:"meting"},{label:"GD Studio",value:"gdstudio"}],p=Qt(),f=z(n[0].id),h=z(!1),u=Zt({}),C=z(""),_=z("meting"),x=z("");function m(d,i){return`${d}::${i}`}function j(d,i){var r;return((r=u[m(d,i)])==null?void 0:r.status)??"idle"}function $(d,i){var r;return((r=u[m(d,i)])==null?void 0:r.error)??""}function y(d,i){var r;return((r=u[m(d,i)])==null?void 0:r.body)??""}function P(d,i){var w;const r=(w=u[m(d,i)])==null?void 0:w.latencyMs;return typeof r=="number"?`${r} ms`:"—"}function T(d,i){var w;const r=(w=u[m(d,i)])==null?void 0:w.checkedAt;return r?new Date(r).toLocaleTimeString():"—"}function B(d){switch(d){case"loading":return"检测中";case"ok":return"正常";case"fail":return"异常";default:return"未检测"}}async function V(d,i){const r=m(d,i.id);u[r]={status:"loading"};const w=performance.now();try{const L=await fetch(i.url,{method:"GET",cache:"no-store",signal:AbortSignal.timeout(8e3)}),H=Math.round(performance.now()-w),te=(await L.text()).slice(0,500);if(!L.ok){u[r]={status:"fail",latencyMs:H,checkedAt:Date.now(),error:`HTTP ${L.status}`,body:te||void 0};return}u[r]={status:"ok",latencyMs:H,checkedAt:Date.now(),body:te||void 0}}catch(L){let H=L instanceof Error?L.message:String(L);/failed to fetch|networkerror|load failed/i.test(H)&&(H="Failed to fetch（多半未放行）"),u[r]={status:"fail",latencyMs:Math.round(performance.now()-w),checkedAt:Date.now(),error:H}}}async function M(d){const i=m(F,d.id);u[i]={status:"loading"};const r=await Ca(d);u[i]={status:r.ok?"ok":"fail",latencyMs:r.latencyMs,checkedAt:Date.now(),error:r.error,body:r.body}}function Y(d){return d==="gdstudio"?"GD Studio":"Meting"}function O(){const d=x.value.trim();d&&(p.addSource({title:C.value,kind:_.value,origin:d}),C.value="",x.value="")}async function S(){h.value=!0;try{if(f.value===F){await Promise.all(p.sources.map(i=>M(i)));return}const d=n.find(i=>i.id===f.value);if(!d)return;await Promise.all(d.cards.map(i=>V(d.id,i)))}finally{h.value=!1}}return be(f,()=>{S()}),Ge(()=>{S()}),(d,i)=>(I(),ea(ra,null,{default:se(()=>[b("div",wa,[b("div",_a,[i[4]||(i[4]=b("h1",{class:"browse-section-title border-l-4 border-blue-500 pl-2"},"监控",-1)),b("button",{type:"button",class:"monitor-btn",disabled:h.value,onClick:S},W(h.value?"检测中…":"刷新当前标签"),9,Ta)]),q(E(va),{value:f.value,"onUpdate:value":i[3]||(i[3]=r=>f.value=r),type:"line",animated:"",class:"monitor-tabs"},{default:se(()=>[(I(),X(de,null,me(n,r=>q(E(We),{key:r.id,name:r.id,tab:r.label},{default:se(()=>[b("div",Pa,[(I(!0),X(de,null,me(r.cards,w=>(I(),X("article",{key:w.id,class:Ae(["monitor-card",`monitor-card--${j(r.id,w.id)}`])},[b("div",Ra,[b("div",null,[b("h2",ka,W(w.title),1),b("p",za,W(w.url),1)]),b("span",$a,W(B(j(r.id,w.id))),1)]),b("dl",La,[b("div",null,[i[5]||(i[5]=b("dt",null,"延迟",-1)),b("dd",null,W(P(r.id,w.id)),1)]),b("div",null,[i[6]||(i[6]=b("dt",null,"最近检测",-1)),b("dd",null,W(T(r.id,w.id)),1)])]),$(r.id,w.id)?(I(),X("p",Ba,W($(r.id,w.id)),1)):y(r.id,w.id)?(I(),X("pre",Aa,W(y(r.id,w.id)),1)):xe("",!0),b("button",{type:"button",class:"monitor-btn monitor-btn--ghost",disabled:h.value,onClick:L=>V(r.id,w)}," 单独检测 ",8,Wa)],2))),128))])]),_:2},1032,["name","tab"])),64)),q(E(We),{name:F,tab:"音乐源"},{default:se(()=>[i[11]||(i[11]=b("p",{class:"monitor-hint"},"取链按卡片从上到下尝试。可改地址、开关、排序；内建源只能停用不能删。",-1)),b("div",Ma,[(I(!0),X(de,null,me(E(p).sources,(r,w)=>(I(),X("article",{key:r.id,class:Ae(["monitor-card",[`monitor-card--${j(F,r.id)}`,{"monitor-card--off":!r.enabled}]])},[b("div",Ea,[b("div",null,[b("h2",ja,W(r.title),1),b("p",Oa,W(Y(r.kind))+" · "+W(E(Oe)(r)),1)]),b("span",Va,W(B(j(F,r.id))),1)]),b("label",Ha,[i[7]||(i[7]=b("span",null,"源地址",-1)),q(E(ye),{value:r.origin,placeholder:E(Oe)(r),size:"small","onUpdate:value":L=>E(p).patchSource(r.id,{origin:L})},null,8,["value","placeholder","onUpdate:value"])]),b("dl",Ga,[b("div",null,[i[8]||(i[8]=b("dt",null,"延迟",-1)),b("dd",null,W(P(F,r.id)),1)]),b("div",null,[i[9]||(i[9]=b("dt",null,"最近检测",-1)),b("dd",null,W(T(F,r.id)),1)])]),$(F,r.id)?(I(),X("p",Da,W($(F,r.id)),1)):y(F,r.id)?(I(),X("pre",Ia,W(y(F,r.id)),1)):xe("",!0),b("div",Fa,[b("button",{type:"button",class:"monitor-btn monitor-btn--ghost",disabled:h.value,onClick:L=>M(r)}," 单独检测 ",8,Na),b("button",{type:"button",class:"monitor-btn monitor-btn--ghost",onClick:L=>E(p).patchSource(r.id,{enabled:!r.enabled})},W(r.enabled?"停用":"启用"),9,Ua),b("button",{type:"button",class:"monitor-btn monitor-btn--ghost",disabled:w===0,onClick:L=>E(p).moveSource(r.id,-1)}," 上移 ",8,Xa),b("button",{type:"button",class:"monitor-btn monitor-btn--ghost",disabled:w===E(p).sources.length-1,onClick:L=>E(p).moveSource(r.id,1)}," 下移 ",8,Ka),r.builtin?xe("",!0):(I(),X("button",{key:0,type:"button",class:"monitor-btn monitor-btn--ghost",onClick:L=>E(p).removeSource(r.id)}," 删除 ",8,Ya))])],2))),128))]),b("form",{class:"monitor-add",onSubmit:ta(O,["prevent"])},[i[10]||(i[10]=b("h3",{class:"monitor-add__title"},"添加源",-1)),q(E(ye),{value:C.value,"onUpdate:value":i[0]||(i[0]=r=>C.value=r),size:"small",placeholder:"名称"},null,8,["value"]),q(E(na),{value:_.value,"onUpdate:value":i[1]||(i[1]=r=>_.value=r),size:"small",options:s},null,8,["value"]),q(E(ye),{value:x.value,"onUpdate:value":i[2]||(i[2]=r=>x.value=r),size:"small",placeholder:"https://example.com/meting"},null,8,["value"]),b("button",{type:"submit",class:"monitor-btn",disabled:!x.value.trim()},"添加",8,qa)],32)]),_:1})]),_:1},8,["value"])])]),_:1}))}}),tr=aa(Ja,[["__scopeId","data-v-df7a86c2"]]);export{tr as default};
