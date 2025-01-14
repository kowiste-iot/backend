<#import "template.ftl" as layout>
<@layout.registrationLayout displayMessage=false; section>
    <#if section = "header">
        ${msg("errorTitle")}
    <#elseif section = "form">
        <div class="error-container">
            <div class="error-card">
                <h2>${msg("errorTitle")}</h2>
                <p class="instruction">${message.summary}</p>
                <#if skipLink??>
                    <#else>
                        <#if client?? && client.baseUrl?has_content>
                            <p><a id="backToApplication" href="${client.baseUrl}">${kcSanitize(msg("backToApplication"))?no_esc}</a></p>
                        </#if>
                </#if>
            </div>
        </div>
    </#if>
</@layout.registrationLayout>