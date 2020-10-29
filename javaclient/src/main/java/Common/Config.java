package Common;

import io.netty.handler.ssl.ApplicationProtocolConfig;

import java.util.Arrays;
import java.util.Collections;
import java.util.List;

public class Config {
    public static int port = 50052;
    public static int portNoTls = 50053;

    public static final String TLS_PROTOCOL = "TLSv1.2";


    private static final String GRPC_EXP_VERSION = "grpc-exp";
    private static final String HTTP2_VERSION = "h2";
    private static final List<String> NEXT_PROTOCOL_VERSIONS =
            Collections.unmodifiableList(Arrays.asList(GRPC_EXP_VERSION, HTTP2_VERSION));

    public static final ApplicationProtocolConfig NPN_AND_ALPN = new ApplicationProtocolConfig(
            ApplicationProtocolConfig.Protocol.NPN_AND_ALPN,
            ApplicationProtocolConfig.SelectorFailureBehavior.NO_ADVERTISE,
            ApplicationProtocolConfig.SelectedListenerFailureBehavior.ACCEPT,
            NEXT_PROTOCOL_VERSIONS);

    public static final List<String> CIPHERS = Collections.unmodifiableList(Arrays
            .asList(
                    // GCM (Galois/Counter Mode) requires JDK 8.
                    "TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384",
                    "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256",
                    "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
                    "TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA",
                    // AES256 requires JCE unlimited strength jurisdiction policy files.
                    "TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA",
                    // GCM (Galois/Counter Mode) requires JDK 8.
                    "TLS_RSA_WITH_AES_128_GCM_SHA256",
                    "TLS_RSA_WITH_AES_128_CBC_SHA",
                    // AES256 requires JCE unlimited strength jurisdiction policy files.
                    "TLS_RSA_WITH_AES_256_CBC_SHA",
                    /** 在原有的默认套件中加入国密算法套件
                     * todo 需要注意的是，如果启用了ECDHE算法，服务器会优先选择ECDHE算法套件，
                     * todo 并且ECDHE算法套件必须要求走双向SSL，客户端和服务器端必须要配置加密证书和签名证书
                     */
                    "TLS_ECDHE_WITH_SM4_SM3",
                    "TLS_ECC_WITH_SM4_SM3"
            ));
}
