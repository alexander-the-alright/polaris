// =============================================================================
// Auth: Alex Celani
// File: gps.go
// Revn: 07-27-2023  4.0
// Func: asynchronously receive messages from another program, print
//
// TODO: figure out how to process gps coordinates
// =============================================================================
// CHANGE LOG
// -----------------------------------------------------------------------------
//*07-13-2023: copied over, removed call to alter()/forward()
//             removed any call to write
//             put read and print in a for loop
//             wrote getFileName() to create a filename based off
//              current time
//             write output to a logfile
//             moved file writing to AFTER checking for fin
//*07-20-2023: added calls to os.Create(), file.Write(), etc
// 07-21-2023: tested filewriting stuff, no avail
// 07-24-2023: added bufio import and buffered file writing
//*07-25-2023: commented bufio stuffs
// 07-25-2023: changed start/end flags to Russian based ot/za
//             imported flag, set up to basic flags
//*07-27-2023: imported rpio, wrote blink()
//             added -r reset flag to turn off green file LED
//             commented rpio changes
//             
// =============================================================================

package main

import ( 
    "bufio"     // NewWriter, Write, Writer.Flush
    // Toggle, rpio.Pin, Close, Open, Output
    "github.com/stianeikeland/go-rpio"
    "flag"      // Bool, String, Parse
    "fmt"       // Fprintf, Println
    "net"       // ResolveTCPAddr, conn.Write,Read,Close
                // ListenTCP, listener.Accept
    "os"        // Args, Stderr, Exit, Create, Write, Close,
                // WriteString
    "strings"   // ToLower, Compare
    "time"      // Now, Format, Sleep, Duration, time.Duration
)


// quick error checking
func check( err error ) {
    // if user checks an error, it better be nil
    if err != nil {
        // put error in string, print in stderr
        fmt.Fprintf( os.Stderr, "Fatal error: %s", err.Error() )
        os.Exit( 2 )        // quit
    }
}


func getFileName() string{
    // capture current time
    t := time.Now()
    // format as YYYYMMDDHHmmss.gps
    return t.Format( "20060102150405" ) + ".gps"
}


// control the LEDs
func blink( p rpio.Pin, d time.Duration, f int ) {
    // unmap memory when done
//    defer rpio.Close()

    if perr == nil {        // only blink if no error
        if d == 0 {         // if duration is 0...
            p.Toggle()      // toggle and leave
        } else {
            // if duration is not 0...
            for i := 0; i < 2 * f; i++ {
                p.Toggle()  // toggle led
                // if this isn't of type time.Duration, shit breaks
                var dur time.Duration = 2 * d
                // duration is for on/off, divide by twice the given
                time.Sleep( time.Second / dur )
            }
        }
    }
}


// function to handle incoming connections
func handleClient( conn net.Conn ) {
    defer conn.Close()  // barring any error, still close connection

    var buf [512]byte   // declare large byte array, store messages

    if *debug {     // XXX debug print
        fmt.Println( "handling" )
    }

    // iterate forever to always read over connection
    for {
        // read n bytes from connection into buffer
        n, err := conn.Read( buf[0:] )
        if err != nil { // erroring on read will simply leave the 
            return      // function so it can start again later
        }

        // essentially a case statement for conn.Read, Compare makes
        // it hard to do that
        // check to see if buf[:n] is start
        if strings.Compare( string( buf[:n] ), "ot" ) == 0 {
            filename := getFileName()   // get name of file
            file, err = os.Create( filename )   // create new file
            check( err )    // make sure file creation worked
            if *debug {     // XXX debug print
                fmt.Println( "created file: ", filename )
            }
            writer = bufio.NewWriter( file )    // create filewriter
        // check to see if buf[:n] is fin
        } else if strings.Compare( string( buf[:n] ), "za" ) == 0 {
            if *debug {     // XXX debug print
                fmt.Println( "file closing" )
            }
            // light up green LED on file close
            blink( gpin, 0, 1 )
            file.Close()    // close file
            // possibly strike exit...
            os.Exit( 2 )        // quit
        // if conn.Read is not start or fin, it's data
        } else {
            // write info to file
            _, err = writer.Write( buf[:n] )            
            check( err )            // make sure write worked
            err = writer.Flush()    // flush write to file
            check( err )            // make sure write worked
            blink( rpin, 4, 2 )     // blink LED in 1/4s, twice
            if *debug {     // XXX debug print
                fmt.Println( "writing" )
            }
        }

        if *debug || *verbose {     // XXX debug print
        // string() only works on byte SLICES so [:] is required
            fmt.Println( "recv: ", string( buf[:n] ) )
        }
    }
}


var file *os.File           // global declaration of file
var writer *bufio.Writer    // global declaration of filewriter
// declare empty flag stuffs
var debug *bool
var service *string
var verbose *bool
var reset *bool
// declare pins and accessories
var rpin rpio.Pin
var gpin rpio.Pin
var perr error


func main() {

    // declare red and green LED pins
    rpin = rpio.Pin( 10 )
    gpin = rpio.Pin( 9 )
    // catch an error, silently. If it fails, I'll just ignore blinks
    perr = rpio.Open()

    defer rpio.Close()      // unmap memory when done with it

    // set pins to output
    rpin.Output()
    gpin.Output()

    // define flag for debug printing
    debug = flag.Bool( "d", false, "enable debug prints" )
    // default :1202 is ip:port, no ip, implies localhost
    // this is the ip and port of the process
    service = flag.String( "s", ":1202", "host:ip of process" )
    // define flag to print received messages as well as write to file
    verbose = flag.Bool( "v", false, "enable print on receipt" )
    // define flag to turn off the file-alert LED
    reset = flag.Bool( "r", false, "turn off receipt LED" )

    flag.Parse()    // parse flags

    if *reset {                 // reset toggle
        blink( gpin, 0, 1 )     // turn off green led
        os.Exit( 1 )            // exeunt
    }

    // "resolve" ip & host according to TCP rules
    tcpAddr, err := net.ResolveTCPAddr( "tcp", *service )
    check( err )    // check error
    if *debug {     // XXX debug print
        fmt.Println( "tcpaddr: ", tcpAddr )
    }

    // bind and "listen" to ip and port, according to tcp rules
    listener, err := net.ListenTCP( "tcp", tcpAddr )
    check( err )    // check error
    if *debug {     // XXX debug print
        fmt.Println( "listener: ", listener )
    }


    // iterate forever
    // TODO i mean i can totally make this more user friendly
    for {
        // accept a connection that makes its way to bound port
        conn, err := listener.Accept()
        if err != nil {     // if connection fails...
            continue        // don't quit program, not fatal error
        }

        if *debug {     // XXX debug print
            addr, _ := conn.RemoteAddr().( *net.TCPAddr )
            fmt.Println( "connected ip: ", addr.IP.String()  )
        }

        // asynchronous function to handle connection to client
        go handleClient( conn )
    }
}

