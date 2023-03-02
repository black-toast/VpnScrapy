package com.test.tly;

import java.util.List;

public class VpnNodeResponse {
    private List<VpnNode> node;

    public List<VpnNode> getNode() {
        return node;
    }

    public void setNode(List<VpnNode> node) {
        this.node = node;
    }

    @Override
    public String toString() {
        return "VpnNodeResponse{" +
                "node=" + node +
                '}';
    }
}
