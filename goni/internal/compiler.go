package internal

import "github.com/lyraproj/goni/goni"

type Compiler struct {
	analyzer *Analyzer
	enc goni.Encoding
	regex *Regex
}

func (cm *Compiler) init(analyzer *Analyzer) {
	cm.analyzer = analyzer
	cm.regex = analyzer.regex
	cm.enc = analyzer.enc
}

var y = `
    protected Compiler(Analyser analyser) {
        this.analyser = analyser;
        this.regex = analyser.regex;
        this.enc = regex.enc;
    }

    final void compile(Node root) {
        prepare();
        compileTree(root);
        finish();
    }

    protected abstract void prepare();
    protected abstract void finish();

    protected abstract void compileAltNode(ListNode node);

    private void compileStringRawNode(StringNode sn) {
        if (sn.length() <= 0) return;
        addCompileString(sn.bytes, sn.p, 1 /*sb*/, sn.length(), false);
    }

    private void compileStringNode(StringNode node) {
        StringNode sn = node;
        if (sn.length() <= 0) return;

        boolean ambig = sn.isAmbig();

        int p, prev;
        p = prev = sn.p;
        int end = sn.end;
        byte[]bytes = sn.bytes;
        int prevLen = enc.length(bytes, p, end);
        p += prevLen;
        int blen = prevLen;

        while (p < end) {
            int len = enc.length(bytes, p, end);
            if (len == prevLen || ambig) {
                blen += len;
            } else {
                addCompileString(bytes, prev, prevLen, blen, ambig);
                prev = p;
                blen = len;
                prevLen = len;
            }
            p += len;
        }
        addCompileString(bytes, prev, prevLen, blen, ambig);
    }

    protected abstract void addCompileString(byte[]bytes, int p, int mbLength, int strLength, boolean ignoreCase);

    protected abstract void compileCClassNode(CClassNode node);
    protected abstract void compileCTypeNode(CTypeNode node);
    protected abstract void compileAnyCharNode();
    protected abstract void compileCallNode(CallNode node);
    protected abstract void compileBackrefNode(BackRefNode node);
    protected abstract void compileCECQuantifierNode(QuantifierNode node);
    protected abstract void compileNonCECQuantifierNode(QuantifierNode node);
    protected abstract void compileOptionNode(EncloseNode node);
    protected abstract void compileEncloseNode(EncloseNode node);
    protected abstract void compileAnchorNode(AnchorNode node);

    protected final void compileTree(Node node) {
        switch (node.getType()) {
        case NodeType.LIST:
            ListNode lin = (ListNode)node;
            do {
                compileTree(lin.value);
            } while ((lin = lin.tail) != null);
            break;

        case NodeType.ALT:
            compileAltNode((ListNode)node);
            break;

        case NodeType.STR:
            StringNode sn = (StringNode)node;
            if (sn.isRaw()) {
                compileStringRawNode(sn);
            } else {
                compileStringNode(sn);
            }
            break;

        case NodeType.CCLASS:
            compileCClassNode((CClassNode)node);
            break;

        case NodeType.CTYPE:
            compileCTypeNode((CTypeNode)node);
            break;

        case NodeType.CANY:
            compileAnyCharNode();
            break;

        case NodeType.BREF:
            compileBackrefNode((BackRefNode)node);
            break;

        case NodeType.CALL:
            if (Config.USE_SUBEXP_CALL) {
                compileCallNode((CallNode)node);
                break;
            } // USE_SUBEXP_CALL
            break;

        case NodeType.QTFR:
            if (Config.USE_CEC) {
                compileCECQuantifierNode((QuantifierNode)node);
            } else {
                compileNonCECQuantifierNode((QuantifierNode)node);
            }
            break;

        case NodeType.ENCLOSE:
            EncloseNode enode = (EncloseNode)node;
            if (enode.isOption()) {
                compileOptionNode(enode);
            } else {
                compileEncloseNode(enode);
            }
            break;

        case NodeType.ANCHOR:
            compileAnchorNode((AnchorNode)node);
            break;

        default:
            // undefined node type
            newInternalException(PARSER_BUG);
        } // switch
    }

    protected final void compileTreeNTimes(Node node, int n) {
        for (int i=0; i<n; i++) compileTree(node);
    }

    protected void newSyntaxException(String message) {
        throw new SyntaxException(message);
    }

    protected void newInternalException(String message) {
        throw new InternalException(message);
    }
}
`
