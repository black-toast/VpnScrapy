package com.test.tly;

import com.google.gson.annotations.SerializedName;

public class VpnNode {
    private String node_name;
    private String node_server;
    private int port;
    private String pass;
    private String node_method;
    private String cipher;
    @SerializedName("ws-path")
    private String ws_path;
    private String tls;

    public String getNode_name() {
        return node_name;
    }

    public void setNode_name(String node_name) {
        this.node_name = node_name;
    }

    public String getNode_server() {
        return node_server;
    }

    public void setNode_server(String node_server) {
        this.node_server = node_server;
    }

    public int getPort() {
        return port;
    }

    public void setPort(int port) {
        this.port = port;
    }

    public String getPass() {
        return pass;
    }

    public void setPass(String pass) {
        this.pass = pass;
    }

    public String getNode_method() {
        return node_method;
    }

    public void setNode_method(String node_method) {
        this.node_method = node_method;
    }

    public String getCipher() {
        return cipher;
    }

    public void setCipher(String cipher) {
        this.cipher = cipher;
    }

    public String getWs_path() {
        return ws_path;
    }

    public void setWs_path(String ws_path) {
        this.ws_path = ws_path;
    }

    public String getTls() {
        return tls;
    }

    public void setTls(String tls) {
        this.tls = tls;
    }

    @Override
    public String toString() {
        return "VpnNode{" +
                "node_name='" + node_name + '\'' +
                ", node_server='" + node_server + '\'' +
                ", port=" + port +
                ", pass='" + pass + '\'' +
                ", node_method='" + node_method + '\'' +
                ", cipher='" + cipher + '\'' +
                ", ws_path='" + ws_path + '\'' +
                ", tls='" + tls + '\'' +
                '}';
    }
}
