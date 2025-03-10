
struct ethernet_h {
	bit<48> dst_addr
	bit<48> src_addr
	bit<16> ether_type
}

struct vlan_tag_h {
	bit<16> pcp_cfi_vid
	bit<16> ether_type
}

struct ipv4_h {
	bit<8> version_ihl
	bit<8> diffserv
	bit<16> total_len
	bit<16> identification
	bit<16> flags_frag_offset
	bit<8> ttl
	bit<8> protocol
	bit<16> hdr_checksum
	bit<32> src_addr
	bit<32> dst_addr
}

struct psa_ingress_output_metadata_t {
	bit<8> class_of_service
	bit<8> clone
	bit<16> clone_session_id
	bit<8> drop
	bit<8> resubmit
	bit<32> multicast_group
	bit<32> egress_port
}

struct psa_egress_output_metadata_t {
	bit<8> clone
	bit<16> clone_session_id
	bit<8> drop
}

struct psa_egress_deparser_input_metadata_t {
	bit<32> egress_port
}

struct send_arg_t {
	bit<32> port
}

header ethernet instanceof ethernet_h
header vlan_tag instanceof vlan_tag_h
header ipv4 instanceof ipv4_h

struct my_ingress_metadata_t {
	bit<32> psa_ingress_input_metadata_ingress_port
	bit<8> psa_ingress_output_metadata_drop
	bit<32> psa_ingress_output_metadata_egress_port
}
metadata instanceof my_ingress_metadata_t

action NoAction args none {
	return
}

action send args instanceof send_arg_t {
	mov m.psa_ingress_output_metadata_egress_port t.port
	return
}

action drop_1 args none {
	mov m.psa_ingress_output_metadata_drop 1
	return
}

table ipv4_host {
	key {
		h.ipv4.dst_addr exact
	}
	actions {
		send
		drop_1
		NoAction
	}
	default_action drop_1 args none const
	size 0x10000
}


apply {
	rx m.psa_ingress_input_metadata_ingress_port
	mov m.psa_ingress_output_metadata_drop 0x0
	extract h.ethernet
	jmpeq INGRESS_PARSER_PARSE_VLAN_TAG h.ethernet.ether_type 0x8100
	jmpeq INGRESS_PARSER_PARSE_IPV4 h.ethernet.ether_type 0x800
	jmp INGRESS_PARSER_ACCEPT
	INGRESS_PARSER_PARSE_VLAN_TAG :	extract h.vlan_tag
	jmpeq INGRESS_PARSER_PARSE_IPV4 h.vlan_tag.ether_type 0x800
	jmp INGRESS_PARSER_ACCEPT
	INGRESS_PARSER_PARSE_IPV4 :	extract h.ipv4
	INGRESS_PARSER_ACCEPT :	table ipv4_host
	jmpneq LABEL_DROP m.psa_ingress_output_metadata_drop 0x0
	emit h.ethernet
	emit h.vlan_tag
	emit h.ipv4
	tx m.psa_ingress_output_metadata_egress_port
	LABEL_DROP :	drop
}
