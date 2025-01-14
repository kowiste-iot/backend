<#import "template.ftl" as layout>
<@layout.registrationLayout displayMessage=!messagesPerField.existsError('username','password') displayInfo=realm.password && realm.registrationAllowed && !registrationDisabled??; section>
    <#if section = "header">
        ${msg("loginAccountTitle")}
    <#elseif section = "form">
        <div class="login-container">
            <div class="login-card">
                <h2>${realm.displayName!''}</h2>
                <form onsubmit="login.disabled = true; return true;" action="${url.loginAction}" method="post">
                    <div class="form-group">
                        <label for="username">${msg("username")}</label>
                        <input tabindex="1" id="username" name="username" value="${(login.username!'')}" type="text" autofocus autocomplete="off"/>
                    </div>

                    <div class="form-group">
                        <label for="password">${msg("password")}</label>
                        <input tabindex="2" id="password" name="password" type="password" autocomplete="off"/>
                    </div>

                    <div class="form-group submit">
                        <input type="hidden" id="id-hidden-input" name="credentialId" <#if auth.selectedCredential??>value="${auth.selectedCredential}"</#if>/>
                        <input tabindex="4" class="submit-button" name="login" id="kc-login" type="submit" value="${msg("doLogIn")}"/>
                    </div>
                </form>

                <#if realm.resetPasswordAllowed>
                    <div class="forgot-password">
                        <a tabindex="5" href="${url.loginResetCredentialsUrl}">${msg("doForgotPassword")}</a>
                    </div>
                </#if>
            </div>
        </div>
    </#if>
</@layout.registrationLayout>