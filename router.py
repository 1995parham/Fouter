"""Custom topology example
"""
__author__='Meti'
from mininet.net import Mininet
from mininet.topo import Topo
from mininet.node import Controller, RemoteController, OVSKernelSwitch, UserSwitch, Node
from mininet.cli import CLI
from mininet.log import setLogLevel
from mininet.link import Link, TCLink

class MyTopo(Topo):
    "Simple router project."

    def build(self, **opts):
        # Add hosts, switches, and a router
        r0 = self.addHost( 'r0', ip='10.0.1.1/24')

        s1 = self.addSwitch( 's1' )
        s2 = self.addSwitch( 's2' )
        s3 = self.addSwitch( 's3' )

        h1 = self.addHost( 'h1', ip="10.0.1.10/24", mac="00:00:00:00:00:01", defaultRoute='via 10.0.1.1')
        h2 = self.addHost( 'h2', ip="10.0.2.10/24", mac="00:00:00:00:00:02", defaultRoute='via 10.0.2.1')
        h3 = self.addHost( 'h3', ip="10.0.3.10/24", mac="00:00:00:00:00:03", defaultRoute='via 10.0.3.1')

        # connect to the controller
        # c0 = net.addController( 'c0', controller=RemoteController, ip='127.0.0.1', port=6633 )

        # Add links
        self.addLink( r0, s1, intfName1='r0-eth1')
        self.addLink( r0, s2, intfName1='r0-eth2')
        self.addLink( r0, s3, intfName1='r0-eth3')

        self.addLink( h1, s1 )
        self.addLink( h2, s2 )
        self.addLink( h3, s3 )

def run():
    topo = MyTopo()

    # start mininet program
    net = Mininet(topo=topo, link=TCLink, autoStaticArp=True, autoSetMacs = True)
    net.start()

    for h in net.hosts:
        print "disable ipv6"
        h.cmd("sysctl -w net.ipv6.conf.all.disable_ipv6=1")
        h.cmd("sysctl -w net.ipv6.conf.default.disable_ipv6=1")
        h.cmd("sysctl -w net.ipv6.conf.lo.disable_ipv6=1")

    # set router ips
    net['r0'].cmd('ifconfig r0-eth2 10.0.2.1 netmask 255.255.255.0')
    net['r0'].cmd('ifconfig r0-eth3 10.0.3.1 netmask 255.255.255.0')

    # run router code
    net['r0'].cmd("./fouter > ./out.log 2>&1 &")

    print "*** Running CLI"
    CLI( net )

    print "*** Stopping network"
    net.stop()

#topos = { 'mytopo': ( lambda: MyTopo() ) }


# if the script is run directly (sudo custom/optical.py):
if __name__ == '__main__':
    setLogLevel( 'info' )
    run()
